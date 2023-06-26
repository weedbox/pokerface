module github.com/weedbox/pokerface

go 1.19

require (
	github.com/google/uuid v1.3.0
	github.com/stretchr/testify v1.8.4
	github.com/weedbox/pokertable v0.0.0-20230621122628-b21e12e170f5
	github.com/weedbox/timebank v0.0.0-20230626195305-39f7a14ece16
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/thoas/go-funk v0.9.3 // indirect
	github.com/weedbox/syncsaga v0.0.0-20230626174843-43af3ced8402 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/weedbox/syncsaga => ../weedbox/syncsaga

//replace github.com/weedbox/timebank => ../weedbox/timebank
