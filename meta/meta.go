package meta

import "encoding/json"

// Meta meta tags on html header. This account for both twitter meta atgs and
// opengraph
type Meta struct {
	OpenGraph *OpenGraph
	Twitter   *Twitter
}

// Map returns a map of meta tags.
func (m *Meta) Map() (map[string]string, error) {
	o := make(map[string]string)
	if m.OpenGraph != nil {
		n, err := toMap(m.OpenGraph)
		if err != nil {
			return nil, err
		}
		for k, v := range n {
			o[k] = v
		}
	}
	if m.Twitter != nil {
		n, err := toMap(m.Twitter)
		if err != nil {
			return nil, err
		}
		for k, v := range n {
			o[k] = v
		}
	}
	return o, nil
}

// OpenGraph meta properties
type OpenGraph struct {
	URL         string `json:"og:url,omitempty"`
	Title       string `json:"og:title,omitempty"`
	Description string `json:"og:description,omitempty"`
	Image       string `json:"og:image,omitempty"`
	Type        string `json:"og:type,omitempty"`
	Locale      string `json:"og:locale,omitempty"`
	Site        string `json:"og:site_name,omitempty"`
}

// Twitter is  a twitter card markup properties.
type Twitter struct {
	Card              string `json:"twitter:card,omitempty"`
	Site              string `json:"twitter:site,omitempty"`
	SiteID            string `json:"twitter:site:id,omitempty"`
	Creator           string `json:"twitter:creator,omitempty"`
	CreatorID         string `json:"twitter:creator:id,omitempty"`
	Description       string `json:"twitter:description,omitempty"`
	Title             string `json:"twitter:title,omitempty"`
	Image             string `json:"twitter:image,omitempty"`
	ImageAlt          string `json:"twitter:image:alt,omitempty"`
	Player            string `json:"twitter:player,omitempty"`
	PlayerWidth       string `json:"twitter:player:width,omitempty"`
	PlayerHeight      string `json:"twitter:player:height,omitempty"`
	PlayerStream      string `json:"twitter:player:stream,omitempty"`
	AppNameIphone     string `json:"twitter:app:name:iphone,omitempty"`
	AppIDIphone       string `json:"twitter:app:id:iphone,omitempty"`
	AppURLIphone      string `json:"twitter:app:url:iphone,omitempty"`
	AppNameIpad       string `json:"twitter:app:name:ipad,omitempty"`
	AppIDIpad         string `json:"twitter:app:id:ipad,omitempty"`
	AppURLIpad        string `json:"twitter:app:url:ipad,omitempty"`
	AppNameGooglePlay string `json:"twitter:app:name:googleplay,omitempty"`
	AppIDGooglePlay   string `json:"twitter:app:id:googleplay,omitempty"`
	AppURLGooglePlay  string `json:"twitter:app:url:googleplay,omitempty"`
}

func toMap(v interface{}) (map[string]string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	o := make(map[string]string)
	err = json.Unmarshal(b, &o)
	if err != nil {
		return nil, err
	}
	return o, nil
}
