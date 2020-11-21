package internal

type Config struct {
	Websites []Website
}

type Website struct {
	Url                  string
	Keywords             []string
	KeywordsNotAppearing []string `yaml:"keywords_not_appearing"`
}
