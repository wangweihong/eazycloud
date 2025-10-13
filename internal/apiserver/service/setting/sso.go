package setting

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"net/url"
	"time"

	"github.com/crewjam/saml"
	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	gx509 "github.com/wangweihong/gotoolbox/pkg/certificate/x509"
	"github.com/wangweihong/gotoolbox/pkg/decoder"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

func (s *settingService) IdentityProviderSAMLMetadataUpset(ctx context.Context, req *iapiserver.IdentityProviderMetadataUpsetRequest) error {
	keyEncode := base64.StdEncoding.EncodeToString(req.KeyEncode)
	certEncode := base64.StdEncoding.EncodeToString(req.CertEncode)

	if _, _, err := gx509.ParseCert(keyEncode, certEncode); err != nil {
		return errors.Errorf("invalid certificate:%v", err)
	}

	data := iapiserver.Setting{
		ObjectMeta: imachinery.ObjectMeta{
			Name: iapiserver.SettingKindSSOSamlIdpMetadata,
			Extend: map[string]any{
				"key":                       keyEncode,
				"cert":                      certEncode,
				"endpoint":                  req.Endpoint,
				"authn_name_id_format":      req.AuthnNameIDFormat,
				"redirect_sso_frontend_url": req.RedirectSSOFrontendURL,
			},
		},
	}

	if err := s.store.Settings().Upsert(ctx, &data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *settingService) IdentityProviderSAMLMetadataGet(ctx context.Context) (*iapiserver.IdentityProviderMetadataGetResponse, error) {
	meta, err := s.store.Settings().GetByName(ctx, iapiserver.SettingKindSSOSamlIdpMetadata)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	xmlData, err := xmlSamlGenerate(meta)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &iapiserver.IdentityProviderMetadataGetResponse{
		Setting:    meta,
		XML:        xmlData,
		DecodeKey:  string(decoder.MustBase64Decode(iapiserver.IdentityProviderMetadata{Setting: meta}.GetKey())),
		DecodeCert: string(decoder.MustBase64Decode(iapiserver.IdentityProviderMetadata{Setting: meta}.GetCert())),
	}, nil
}

func (s *settingService) ServiceProviderSAMLMetadataUpset(ctx context.Context, req *iapiserver.ServiceProviderMetadataUpsetRequest) error {
	// data := iapiserver.Setting{
	// 	ObjectMeta: imachinery.ObjectMeta{
	// 		Name: iapiserver.SettingKindSSOSamlSpMetadata,
	// 		Extend: map[string]any{
	// 			"key":                       req.KeyEncode,
	// 			"cert":                      req.CertEncode,
	// 			"endpoint":                  req.Endpoint,
	// 			"authn_name_id_format":      req.AuthnNameIDFormat,
	// 			"redirect_sso_frontend_url": req.RedirectSSOFrontendURL,
	// 		},
	// 	},
	// }

	// if err := s.store.Settings().Upsert(ctx, &data); err != nil {
	// 	return errors.WithStack(err)
	// }
	return nil
}

func (s *settingService) ServiceProviderSAMLMetadataGet(ctx context.Context) (*iapiserver.ServiceProviderMetadataGetResponse, error) {
	// meta, err := s.store.Settings().GetByName(ctx, iapiserver.SettingKindSSOSamlSpMetadata)
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }
	// xmlData, err := xmlSamlGenerate(meta)
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }
	// return &iapiserver.ServiceProviderMetadataGetResponse{
	// 	Setting:    meta,
	// 	XML:        xmlData,
	// 	DecodeKey:  string(decoder.MustBase64Decode(iapiserver.IdentityProviderMetadata{Setting: meta}.GetKey())),
	// 	DecodeCert: string(decoder.MustBase64Decode(iapiserver.IdentityProviderMetadata{Setting: meta}.GetCert())),
	// }, nil
	return nil, nil
}

type IdentityProvider struct {
	Key             crypto.PrivateKey
	Signer          crypto.Signer
	Certificate     *x509.Certificate
	Intermediates   []*x509.Certificate
	RootURL         url.URL
	MetadataURL     url.URL
	SSOURL          url.URL
	LogoutURL       url.URL
	SignatureMethod string
	ValidDuration   *time.Duration
}

func xmlSamlGenerate(meta *iapiserver.Setting) (string, error) {
	_, ed, err := samlEntityDescriptorGenerate(meta)
	if err != nil {
		return "", err
	}

	buf, err := xml.MarshalIndent(ed, "", "  ")
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func samlEntityDescriptorGenerate(meta *iapiserver.Setting) (*IdentityProvider, *saml.EntityDescriptor, error) {
	idpMeta := iapiserver.IdentityProviderMetadata{Setting: meta}

	keyPair, Leaf, err := gx509.ParseCert(idpMeta.GetCert(), idpMeta.GetKey())
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	rootURL, err := url.Parse(idpMeta.GetEndpoint())
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	ssoURL := rootURL.ResolveReference(&url.URL{Path: iapiserver.UriSSO})
	//	sloURL := rootURL.ResolveReference(&url.URL{Path: slo_uri})
	validDuration := iapiserver.DefaultCertificateExpireDuration
	certStr := base64.StdEncoding.EncodeToString(Leaf.Raw)

	ed := &saml.EntityDescriptor{
		EntityID:      rootURL.String(),
		ValidUntil:    time.Now().Add(validDuration),
		CacheDuration: validDuration,
		IDPSSODescriptors: []saml.IDPSSODescriptor{
			{
				SSODescriptor: saml.SSODescriptor{
					RoleDescriptor: saml.RoleDescriptor{
						ProtocolSupportEnumeration: "urn:oasis:names:tc:SAML:2.0:protocol",
						KeyDescriptors: []saml.KeyDescriptor{
							{
								Use: "signing",
								KeyInfo: saml.KeyInfo{
									X509Data: saml.X509Data{
										X509Certificates: []saml.X509Certificate{
											{Data: certStr},
										},
									},
								},
							},
							{
								Use: "encryption",
								KeyInfo: saml.KeyInfo{
									X509Data: saml.X509Data{
										X509Certificates: []saml.X509Certificate{
											{Data: certStr},
										},
									},
								},
								EncryptionMethods: []saml.EncryptionMethod{
									{Algorithm: "http://www.w3.org/2001/04/xmlenc#aes128-cbc"},
									{Algorithm: "http://www.w3.org/2001/04/xmlenc#aes192-cbc"},
									{Algorithm: "http://www.w3.org/2001/04/xmlenc#aes256-cbc"},
									{Algorithm: "http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p"},
								},
							},
						},
					},
					NameIDFormats: []saml.NameIDFormat{saml.NameIDFormat("urn:oasis:names:tc:SAML:2.0:nameid-format:transient")},
				},
				SingleSignOnServices: []saml.Endpoint{
					{
						Binding:  saml.HTTPRedirectBinding,
						Location: ssoURL.String(),
					},
					{
						Binding:  saml.HTTPPostBinding,
						Location: ssoURL.String(),
					},
				},
			},
		},
	}

	//ed.IDPSSODescriptors[0].SSODescriptor.SingleLogoutServices = []saml.Endpoint{
	//	{
	//		Binding:  saml.HTTPRedirectBinding,
	//		Location: sloURL.String(),
	//	},
	//}

	idp := &IdentityProvider{
		Key:             keyPair.PrivateKey,
		Signer:          nil,
		Certificate:     keyPair.Leaf,
		Intermediates:   nil,
		RootURL:         *rootURL,
		MetadataURL:     *rootURL.ResolveReference(&url.URL{Path: iapiserver.UriSSOMetadata}),
		SSOURL:          *rootURL.ResolveReference(&url.URL{Path: iapiserver.UriSSO}),
		LogoutURL:       *rootURL.ResolveReference(&url.URL{Path: iapiserver.UriSLO}),
		SignatureMethod: "",
	}

	return idp, ed, nil
}
