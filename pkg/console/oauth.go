package console

import (
	"fmt"
	oauthv1 "github.com/openshift/api/oauth/v1"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/openshift/console-operator/pkg/crypto"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


func newOauthConfigSecret(randomSecret string) *corev1.Secret {
	meta := sharedMeta()
	meta.Name = consoleOauthConfigName
	oauthConfigSecret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind: "Secret",
		},
		ObjectMeta: meta,
		StringData: map[string]string{
			"clientSecret": randomSecret,
		},
	}
	return oauthConfigSecret
}

// the oauth client can be created after the route, once we have a hostname
// - will create a client secret
//   - reference by configmap/deployment
func newConsoleOauthClient(rt *routev1.Route) (*oauthv1.OAuthClient, *corev1.Secret) {
	randomBits := crypto.RandomBitsString(256)
	oauthConfigSecret := newOauthConfigSecret(randomBits)
	host := rt.Spec.Host
	oauthclient := &oauthv1.OAuthClient{
		TypeMeta: metav1.TypeMeta{
			APIVersion: oauthv1.GroupVersion.String(),
			Kind: "OAuthClient",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            "console-oauth-client",
			// logically we own,
			// but namespaced resources cannot own
			// cluster scoped resources
			// OwnerReferences:            nil,
		},
		Secret: randomBits,
		// TODO: we need to fill this in from our Route, whenever
		// it gets a .Spec.Host
		//redirectURIs:
		//- http://localhost:9000/auth/callback
		RedirectURIs:                        []string{
			host,
		},
	}

	// TODO: its a little weird this function returns two objects,
	// but that might be because I am new to golang
	return oauthclient, oauthConfigSecret
}

func UpdateOauthClient(rt *routev1.Route) {
	oauthClient, _ := newConsoleOauthClient(rt)
	fmt.Println("new oauth client with host:")
	logYaml(oauthClient)
	sdk.Update(oauthClient)
}