package constants

const ConfigFileName = ".nextjs_routing_helper.json"

type Config struct {
	Router              RouterType         `json:"router"`
	Language            LanguageType       `json:"language"`
	ComponentStyle      ComponentStyleType `json:"componentStyle"`
	SrcFolder           bool               `json:"srcFolder"`
	PageComponentSuffix string             `json:"pageComponentSuffix"`
}
