package i18n

var _ TranslateInterface = (*Translator)(nil)
var _ TranslateServiceInterface = (*TranslatorService)(nil)

type TranslateServiceInterface interface {
	GetLocal() string
	F(format string, args ...interface{}) string
	T(key string) string
	NStr(m interface{}, args ...interface{}) string
}

type TranslateInterface interface {
	SetLocal(local string) TranslateServiceInterface
	TranslateServiceInterface
}

var instanceTranslator *Translator

func Translotor() *Translator {

	if instanceTranslator == nil {
		instanceTranslator = &Translator{}
		instanceTranslator.init()
	}

	return instanceTranslator

}

type Translator struct {
}

type TranslatorService struct {
	local string
}

func (t *Translator) GetLocalTranslator(format string, args ...interface{}) string {
	return ""
}

func (t *TranslatorService) F(format string, args ...interface{}) string {
	return ""
}
func (t *TranslatorService) T(key string) string {
	return instanceTranslator.T(key)
}
func (t *TranslatorService) NStr(m interface{}, args ...interface{}) string {
	return ""
}

func (t *TranslatorService) GetLocal() string {
	return t.local
}

func (t *Translator) init() {

}

func (t *Translator) SetLocal(local string) TranslateServiceInterface {

	_ = Translotor()

	return &TranslatorService{
		local,
	}

}

func (t *Translator) GetLocal() string {
	return ""
}

func (t *Translator) F(format string, args ...interface{}) string {
	return ""
}
func (t *Translator) T(key string) string {
	return ""
}
func (t *Translator) NStr(m interface{}, args ...interface{}) string {
	return ""
}
