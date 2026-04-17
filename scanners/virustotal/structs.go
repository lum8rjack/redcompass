package virustotal

import "time"

type VT_Response struct {
	Data struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Attributes struct {
			WhoisDate          int    `json:"whois_date"`
			Tld                string `json:"tld"`
			LastDNSRecordsDate int    `json:"last_dns_records_date"`
			LastAnalysisDate   int    `json:"last_analysis_date"`
			Categories         struct {
			} `json:"categories"`
			LastUpdateDate int    `json:"last_update_date"`
			Registrar      string `json:"registrar"`
			TotalVotes     struct {
				Harmless  int `json:"harmless"`
				Malicious int `json:"malicious"`
			} `json:"total_votes"`
			LastHTTPSCertificate struct {
				Size      int `json:"size"`
				PublicKey struct {
					Rsa struct {
						KeySize  int    `json:"key_size"`
						Modulus  string `json:"modulus"`
						Exponent string `json:"exponent"`
					} `json:"rsa"`
					Algorithm string `json:"algorithm"`
				} `json:"public_key"`
				ThumbprintSha256 string `json:"thumbprint_sha256"`
				Tags             []any  `json:"tags"`
				CertSignature    struct {
					SignatureAlgorithm string `json:"signature_algorithm"`
					Signature          string `json:"signature"`
				} `json:"cert_signature"`
				Validity struct {
					NotAfter  string `json:"not_after"`
					NotBefore string `json:"not_before"`
				} `json:"validity"`
				Version    string `json:"version"`
				Extensions struct {
					CertificatePolicies    []string `json:"certificate_policies"`
					ExtendedKeyUsage       []string `json:"extended_key_usage"`
					Tags                   []any    `json:"tags"`
					SubjectAlternativeName []string `json:"subject_alternative_name"`
					AuthorityKeyIdentifier struct {
						Keyid string `json:"keyid"`
					} `json:"authority_key_identifier"`
					CaInformationAccess struct {
						CAIssuers string `json:"CA Issuers"`
						Ocsp      string `json:"OCSP"`
					} `json:"ca_information_access"`
					SubjectKeyIdentifier  string   `json:"subject_key_identifier"`
					CrlDistributionPoints []string `json:"crl_distribution_points"`
					KeyUsage              []string `json:"key_usage"`
					One3614111129242      string   `json:"1.3.6.1.4.1.11129.2.4.2"`
					Ca                    bool     `json:"CA"`
				} `json:"extensions"`
				SignatureAlgorithm string `json:"signature_algorithm"`
				SerialNumber       string `json:"serial_number"`
				Thumbprint         string `json:"thumbprint"`
				Issuer             struct {
					C  string `json:"C"`
					Cn string `json:"CN"`
					O  string `json:"O"`
				} `json:"issuer"`
				Subject struct {
					Cn string `json:"CN"`
					C  string `json:"C"`
					L  string `json:"L"`
					O  string `json:"O"`
					St string `json:"ST"`
				} `json:"subject"`
			} `json:"last_https_certificate"`
			Jarm                 string   `json:"jarm"`
			Tags                 []string `json:"tags"`
			LastModificationDate int      `json:"last_modification_date"`
			Rdap                 struct {
				LdhName     string `json:"ldh_name"`
				UnicodeName string `json:"unicode_name"`
				Events      []struct {
					EventAction string    `json:"event_action"`
					EventDate   time.Time `json:"event_date"`
					EventActor  string    `json:"event_actor"`
					Links       []any     `json:"links"`
				} `json:"events"`
				Links []struct {
					Href     string `json:"href"`
					Rel      string `json:"rel"`
					Type     string `json:"type"`
					Value    string `json:"value"`
					Title    string `json:"title"`
					Media    string `json:"media"`
					HrefLang []any  `json:"href_lang"`
				} `json:"links"`
				Notices []struct {
					Title       string   `json:"title"`
					Description []string `json:"description"`
					Links       []struct {
						Href     string `json:"href"`
						Rel      string `json:"rel"`
						Type     string `json:"type"`
						Value    string `json:"value"`
						Title    string `json:"title"`
						Media    string `json:"media"`
						HrefLang []any  `json:"href_lang"`
					} `json:"links"`
					Type string `json:"type"`
				} `json:"notices"`
				Nameservers []struct {
					LdhName         string `json:"ldh_name"`
					UnicodeName     string `json:"unicode_name"`
					ObjectClassName string `json:"object_class_name"`
					Handle          string `json:"handle"`
					Events          []any  `json:"events"`
					Links           []any  `json:"links"`
					Notices         []any  `json:"notices"`
					Status          []any  `json:"status"`
					Lang            string `json:"lang"`
					Port43          string `json:"port43"`
					Entities        []any  `json:"entities"`
					Remarks         []any  `json:"remarks"`
				} `json:"nameservers"`
				RdapConformance []string `json:"rdap_conformance"`
				Entities        []struct {
					ObjectClassName string `json:"object_class_name"`
					Handle          string `json:"handle"`
					VcardArray      []struct {
						Name       string   `json:"name"`
						Type       string   `json:"type"`
						Values     []string `json:"values"`
						Parameters struct {
						} `json:"parameters,omitempty"`
						Parameters0 struct {
							Type []string `json:"type"`
						} `json:"parameters,omitempty"`
					} `json:"vcard_array"`
					Roles     []string `json:"roles"`
					PublicIds []struct {
						Type       string `json:"type"`
						Identifier string `json:"identifier"`
					} `json:"public_ids"`
					Entities []struct {
						ObjectClassName string `json:"object_class_name"`
						VcardArray      []struct {
							Name       string   `json:"name"`
							Type       string   `json:"type"`
							Values     []string `json:"values"`
							Parameters struct {
							} `json:"parameters,omitempty"`
							Parameters0 struct {
								Type []string `json:"type"`
							} `json:"parameters,omitempty"`
						} `json:"vcard_array"`
						Roles           []string `json:"roles"`
						Handle          string   `json:"handle"`
						PublicIds       []any    `json:"public_ids"`
						Entities        []any    `json:"entities"`
						Remarks         []any    `json:"remarks"`
						Links           []any    `json:"links"`
						Events          []any    `json:"events"`
						AsEventActor    []any    `json:"as_event_actor"`
						Status          []any    `json:"status"`
						Port43          string   `json:"port43"`
						Networks        []any    `json:"networks"`
						Autnums         []any    `json:"autnums"`
						URL             string   `json:"url"`
						Lang            string   `json:"lang"`
						RdapConformance []any    `json:"rdap_conformance"`
					} `json:"entities"`
					Links []struct {
						Href     string `json:"href"`
						Rel      string `json:"rel"`
						Value    string `json:"value"`
						Type     string `json:"type"`
						Title    string `json:"title"`
						Media    string `json:"media"`
						HrefLang []any  `json:"href_lang"`
					} `json:"links"`
					Remarks         []any  `json:"remarks"`
					Events          []any  `json:"events"`
					AsEventActor    []any  `json:"as_event_actor"`
					Status          []any  `json:"status"`
					Port43          string `json:"port43"`
					Networks        []any  `json:"networks"`
					Autnums         []any  `json:"autnums"`
					URL             string `json:"url"`
					Lang            string `json:"lang"`
					RdapConformance []any  `json:"rdap_conformance"`
				} `json:"entities"`
				ObjectClassName string   `json:"object_class_name"`
				Status          []string `json:"status"`
				SecureDNS       struct {
					DelegationSigned bool  `json:"delegation_signed"`
					ZoneSigned       bool  `json:"zone_signed"`
					MaxSigLife       int   `json:"max_sig_life"`
					DsData           []any `json:"ds_data"`
					KeyData          []any `json:"key_data"`
				} `json:"secure_dns"`
				Port43     string `json:"port43"`
				Lang       string `json:"lang"`
				Handle     string `json:"handle"`
				Punycode   string `json:"punycode"`
				Type       string `json:"type"`
				SwitchName string `json:"switch_name"`
				PublicIds  []any  `json:"public_ids"`
				Remarks    []any  `json:"remarks"`
				Nask0State string `json:"nask0_state"`
				Variants   []any  `json:"variants"`
			} `json:"rdap"`
			LastAnalysisStats struct {
				Malicious  int `json:"malicious"`
				Suspicious int `json:"suspicious"`
				Undetected int `json:"undetected"`
				Harmless   int `json:"harmless"`
				Timeout    int `json:"timeout"`
			} `json:"last_analysis_stats"`
			ExpirationDate           int    `json:"expiration_date"`
			Whois                    string `json:"whois"`
			LastHTTPSCertificateDate int    `json:"last_https_certificate_date"`
			LastAnalysisResults      struct {
				Acronis struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Acronis"`
				ZeroXSIF33D struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"0xSI_f33d"`
				Abusix struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Abusix"`
				ADMINUSLabs struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"ADMINUSLabs"`
				Axur struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Axur"`
				ChainPatrol struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"ChainPatrol"`
				CriminalIP struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Criminal IP"`
				AILabsMONITORAPP struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"AILabs (MONITORAPP)"`
				AlienVault struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"AlienVault"`
				AlphaMountainAi struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"alphaMountain.ai"`
				AlphaSOC struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"AlphaSOC"`
				AntiyAVL struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Antiy-AVL"`
				ArcSightThreatIntelligence struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"ArcSight Threat Intelligence"`
				AutoShun struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"AutoShun"`
				BenkowCc struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"benkow.cc"`
				BforeAiPreCrime struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Bfore.Ai PreCrime"`
				BitDefender struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"BitDefender"`
				Bkav struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Bkav"`
				Blueliv struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Blueliv"`
				Certego struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Certego"`
				ChongLuaDao struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Chong Lua Dao"`
				CINSArmy struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"CINS Army"`
				Cluster25 struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Cluster25"`
				Crdf struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"CRDF"`
				CSISSecurityGroup struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"CSIS Security Group"`
				SnortIPSampleList struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Snort IP sample list"`
				CMCThreatIntelligence struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"CMC Threat Intelligence"`
				CTXAI struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"CTX AI"`
				Cyan struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Cyan"`
				Cyble struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Cyble"`
				CyRadar struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"CyRadar"`
				DNS8 struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"DNS8"`
				DrWeb struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Dr.Web"`
				Ermes struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Ermes"`
				Eset struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"ESET"`
				ESTsecurity struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"ESTsecurity"`
				EmergingThreats struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"EmergingThreats"`
				Emsisoft struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Emsisoft"`
				ForcepointThreatSeeker struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Forcepoint ThreatSeeker"`
				Fortinet struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Fortinet"`
				GData struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"G-Data"`
				GCPAbuseIntelligence struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"GCP Abuse Intelligence"`
				GoogleSafebrowsing struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Google Safebrowsing"`
				GreenSnow struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"GreenSnow"`
				GreyNoise struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"GreyNoise"`
				Gridinsoft struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Gridinsoft"`
				Guardpot struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Guardpot"`
				HeimdalSecurity struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Heimdal Security"`
				HuntIoIntelligence struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Hunt.io Intelligence"`
				IPsum struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"IPsum"`
				JuniperNetworks struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Juniper Networks"`
				Kaspersky struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Kaspersky"`
				LevelBlue struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"LevelBlue"`
				Lionic struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Lionic"`
				Lumu struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Lumu"`
				MalwarePatrol struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"MalwarePatrol"`
				MalwareURL struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"MalwareURL"`
				Malwared struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Malwared"`
				Mimecast struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Mimecast"`
				Netcraft struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Netcraft"`
				OpenPhish struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"OpenPhish"`
				PhishingDatabase struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Phishing Database"`
				PhishFort struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"PhishFort"`
				PhishLabs struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"PhishLabs"`
				Phishtank struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Phishtank"`
				Prebytes struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"PREBYTES"`
				PrecisionSec struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"PrecisionSec"`
				QuickHeal struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Quick Heal"`
				Quttera struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Quttera"`
				SafeToOpen struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"SafeToOpen"`
				SansecEComscan struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Sansec eComscan"`
				Scantitan struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Scantitan"`
				SCUMWAREOrg struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"SCUMWARE.org"`
				Seclookup struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Seclookup"`
				SecureBrain struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"SecureBrain"`
				SOCRadar struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"SOCRadar"`
				Sophos struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Sophos"`
				Spam404 struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Spam404"`
				StopForumSpam struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"StopForumSpam"`
				SucuriSiteCheck struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Sucuri SiteCheck"`
				ThreatHive struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"ThreatHive"`
				URLhaus struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"URLhaus"`
				URLQuery struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"URLQuery"`
				ViettelThreatIntelligence struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Viettel Threat Intelligence"`
				Vipre struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"VIPRE"`
				VXVault struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"VX Vault"`
				ViriBack struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"ViriBack"`
				Webroot struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Webroot"`
				YandexSafebrowsing struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Yandex Safebrowsing"`
				ZeroCERT struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"ZeroCERT"`
				DesenmascaraMe struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"desenmascara.me"`
				Securolytics struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"securolytics"`
				XcitiumVerdictCloud struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"Xcitium Verdict Cloud"`
				ZeroFox struct {
					Method     string `json:"method"`
					EngineName string `json:"engine_name"`
					Category   string `json:"category"`
					Result     string `json:"result"`
				} `json:"ZeroFox"`
			} `json:"last_analysis_results"`
			LastDNSRecords []struct {
				Type     string `json:"type"`
				TTL      int    `json:"ttl"`
				Priority int    `json:"priority,omitempty"`
				Value    string `json:"value"`
				Rname    string `json:"rname,omitempty"`
				Serial   int64  `json:"serial,omitempty"`
				Refresh  int    `json:"refresh,omitempty"`
				Retry    int    `json:"retry,omitempty"`
				Expire   int    `json:"expire,omitempty"`
				Minimum  int    `json:"minimum,omitempty"`
			} `json:"last_dns_records"`
			PopularityRanks struct {
			} `json:"popularity_ranks"`
			Reputation   int `json:"reputation"`
			CreationDate int `json:"creation_date"`
		} `json:"attributes"`
	} `json:"data"`
}

type VT_Database struct {
	Domain              string                               `json:"domain"`
	VotesHarmless       int                                  `json:"votes_harmless"`
	VotesMalicious      int                                  `json:"votes_malicious"`
	Malicious           int                                  `json:"malicious"`
	Suspicious          int                                  `json:"suspicious"`
	Undetected          int                                  `json:"undetected"`
	Harmless            int                                  `json:"harmless"`
	Timeout             int                                  `json:"timeout"`
	LastAnalysisResults map[string]VT_Analysis_Engine_Result `json:"last_analysis_results"`
}

type VT_Analysis_Engine_Result struct {
	Method     string `json:"method"`
	EngineName string `json:"engine_name"`
	Category   string `json:"category"`
	Result     string `json:"result"`
}

// VT_QuotaPathCounts maps API path patterns (e.g. "/api/v3/(domains)") to request counts.
type VT_QuotaPathCounts map[string]int

// VT_QuotaResponse is the JSON envelope for VirusTotal usage / daily quota payloads.
type VT_QuotaResponse struct {
	Data VT_QuotaData `json:"data"`
}

// VT_QuotaData holds cumulative totals and per-day usage keyed by calendar date (YYYY-MM-DD).
type VT_QuotaData struct {
	Total VT_QuotaPathCounts            `json:"total"`
	Daily map[string]VT_QuotaPathCounts `json:"daily"`
}
