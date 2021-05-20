package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssmincidents"
	"github.com/fujiwara/ridge"
	"github.com/pkg/errors"
)

var responsePlanArn string

func main() {
	responsePlanArn = os.Getenv("RESPONSE_PLAN_ARN")
	log.Println("[info] checking self IP address")
	resp, err := http.Get("http://checkip.amazonaws.com/")
	if err != nil {
		log.Println("[warn]", err)
	} else {
		io.Copy(os.Stderr, resp.Body)
		resp.Body.Close()
	}
	var mux = http.NewServeMux()
	mux.HandleFunc("/webhook", handleWebhook)
	ridge.Run(":8000", "/", mux)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if t := r.Header.Get("content-type"); !strings.HasPrefix(t, "application/json") {
		errorResponse(w, http.StatusBadRequest, errors.Errorf("invalid content-type %s", t))
		return
	}
	var hook MackerelWebhook
	err := json.NewDecoder(r.Body).Decode(&hook)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
	}
	if s := hook.Alert.Status; s != "critical" {
		log.Printf("[info] alert status is %s. not a critical. ignored.", s)
		return
	}

	log.Println("[info] start incident:", hook.IncidentTitle(), hook.Alert.URL)
	sess := session.Must(session.NewSession())
	svc := ssmincidents.New(sess)
	out, err := svc.StartIncident(&ssmincidents.StartIncidentInput{
		Title:           aws.String(hook.IncidentTitle()),
		ResponsePlanArn: aws.String(responsePlanArn),
		RelatedItems: []*ssmincidents.RelatedItem{
			{
				Title: aws.String("Mackerel"),
				Identifier: &ssmincidents.ItemIdentifier{
					Type: aws.String("OTHER"),
					Value: &ssmincidents.ItemValue{
						Url: &hook.Alert.URL,
					},
				},
			},
		},
	})
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	log.Printf("[info] incident record arn: %s", *out.IncidentRecordArn)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func errorResponse(w http.ResponseWriter, code int, err error) {
	log.Printf("[error] %d %s", code, err)
	w.WriteHeader(code)
}
