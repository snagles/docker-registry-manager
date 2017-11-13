package schema1

import (
	"crypto/x509"

<<<<<<< 8cfda1af4ac8d899bc4d22049fdc575b5d848a69
<<<<<<< e92f10ae507e72db9910b4de55c170515c11ed98
=======
	"github.com/sirupsen/logrus"
>>>>>>> Rename sirupsen/logrus to lower case
=======
>>>>>>> Update dependencies, switch to dep
	"github.com/docker/libtrust"
	"github.com/sirupsen/logrus"
)

// Verify verifies the signature of the signed manifest returning the public
// keys used during signing.
func Verify(sm *SignedManifest) ([]libtrust.PublicKey, error) {
	js, err := libtrust.ParsePrettySignature(sm.all, "signatures")
	if err != nil {
		logrus.WithField("err", err).Debugf("(*SignedManifest).Verify")
		return nil, err
	}

	return js.Verify()
}

// VerifyChains verifies the signature of the signed manifest against the
// certificate pool returning the list of verified chains. Signatures without
// an x509 chain are not checked.
func VerifyChains(sm *SignedManifest, ca *x509.CertPool) ([][]*x509.Certificate, error) {
	js, err := libtrust.ParsePrettySignature(sm.all, "signatures")
	if err != nil {
		return nil, err
	}

	return js.VerifyChains(ca)
}
