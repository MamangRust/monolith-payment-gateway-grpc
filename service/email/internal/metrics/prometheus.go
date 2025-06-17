package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	EmailSent = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "email_sent_total",
		Help: "Total emails sent successfully",
	})

	EmailFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "email_failed_total",
		Help: "Total emails failed",
	})

	EmailRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "email_requests_total",
			Help: "Total email send attempts",
		},
		[]string{"result"},
	)
)

func Register() {
	prometheus.MustRegister(EmailSent, EmailFailed, EmailRequests)
}
