package prometheusrule

import (
	"os"
	"regexp"

	"github.com/openshift/managed-cluster-validating-webhooks/pkg/webhooks/utils"
	admissionv1 "k8s.io/api/admission/v1"
	admissionregv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	admissionctl "sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	WebhookName               string = "promrule-validation"
	privilegedServiceAccounts string = `^system:serviceaccounts:(kube.*|openshift.*|default|redhat.*|osde2e-[a-z0-9]{5})`
	SrepNamespace             string = `(openshift-*|kube-*|redhat-*)`
	docString                 string = `Managed OpenShift Customers may not create promrule in namespaces managed by SRE .`
	clusterAdminGroup         string = "cluster-admins"
)

var (
	timeout                     int32 = 2
	log                               = logf.Log.WithName(WebhookName)
	clusterAdminUsers                 = []string{"kube:admin", "system:admin", "backplane-cluster-admin"}
	sreAdminGroups                    = []string{"system:serviceaccounts:openshift-backplane-srep"}
	privilegedServiceAccountsRe       = regexp.MustCompile(privilegedServiceAccounts)
	privilegedNamespaceRe             = regexp.MustCompile(SrepNamespace)

	scope = admissionregv1.ClusterScope
	rules = []admissionregv1.RuleWithOperations{
		{
			Operations: []admissionregv1.OperationType{"CREATE"},
			Rule: admissionregv1.Rule{
				APIGroups:   []string{""},
				APIVersions: []string{"*"},
				Resources:   []string{"prometheusrule"},
				Scope:       &scope,
			},
		},
	}
)

// PrometheusruleWebhook validates a promrulechange
type PrometheusruleWebhook struct {
	s runtime.Scheme
}

// NewWebhook creates the new webhook
func NewWebhook() *PrometheusruleWebhook {
	scheme := runtime.NewScheme()
	err := admissionv1.AddToScheme(scheme)
	if err != nil {
		log.Error(err, "Fail adding admissionsv1 scheme to PrometheusruleWebhook")
		os.Exit(1)
	}
	err = corev1.AddToScheme(scheme)
	if err != nil {
		log.Error(err, "Fail adding corev1 scheme to PrometheusruleWebhook")
		os.Exit(1)
	}

	return &PrometheusruleWebhook{
		s: *scheme,
	}
}

// Authorized implements Webhook interface
func (s *PrometheusruleWebhook) Authorized(request admissionctl.Request) admissionctl.Response {
	return s.authorized(request)
}

func (s *PrometheusruleWebhook) authorized(request admissionctl.Request) admissionctl.Response {

}

// GetURI implements Webhook interface
func (s *PrometheusruleWebhook) GetURI() string {
	return "/" + WebhookName
}

// Validate implements Webhook interface
func (s *PrometheusruleWebhook) Validate(request admissionctl.Request) bool {

}

// Name implements Webhook interface
func (s *PrometheusruleWebhook) Name() string {
	return WebhookName
}

// FailurePolicy implements Webhook interface
func (s *PrometheusruleWebhook) FailurePolicy() admissionregv1.FailurePolicyType {
	return admissionregv1.Ignore
}

// MatchPolicy implements Webhook interface
func (s *PrometheusruleWebhook) MatchPolicy() admissionregv1.MatchPolicyType {
	return admissionregv1.Equivalent
}

// Rules implements Webhook interface
func (s *PrometheusruleWebhook) Rules() []admissionregv1.RuleWithOperations {
	return rules
}

// ObjectSelector implements Webhook interface
func (s *PrometheusruleWebhook) ObjectSelector() *metav1.LabelSelector {
	return nil
}

// SideEffects implements Webhook interface
func (s *PrometheusruleWebhook) SideEffects() admissionregv1.SideEffectClass {
	return admissionregv1.SideEffectClassNone
}

// TimeoutSeconds implements Webhook interface
func (s *PrometheusruleWebhook) TimeoutSeconds() int32 {
	return timeout
}

// Doc implements Webhook interface
func (s *PrometheusruleWebhook) Doc() string {

}

// SyncSetLabelSelector returns the label selector to use in the SyncSet.
// Return utils.DefaultLabelSelector() to stick with the default
func (s *PrometheusruleWebhook) SyncSetLabelSelector() metav1.LabelSelector {
	return utils.DefaultLabelSelector()
}
