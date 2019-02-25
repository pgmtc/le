package common

//go:generate mockgen -destination=./mocks/mock_marshaller.go -package=mocks github.com/pgmtc/le/pkg/common Marshaller

type Marshaller interface {
	Marshall(data interface{}, fileName string) (resultErr error)
	Unmarshall(fileName string, out interface{}) (resultErr error)
}

type YamlMarshaller struct{}

func (YamlMarshaller) Marshall(data interface{}, fileName string) (resultErr error) {
	return YamlMarshall(data, fileName)
}

func (YamlMarshaller) Unmarshall(fileName string, out interface{}) (resultErr error) {
	return YamlUnmarshall(fileName, out)
}
