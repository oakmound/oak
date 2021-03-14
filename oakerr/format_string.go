package oakerr

import "fmt"

type errCode int

const (
	codeNotFound errCode = iota
	codeExistingElement
	codeExistingElementOverwritten
	codeInsufficientInputs
	codeUnsupportedFormat
	codeNilInput
	codeIndivisibleInput
	codeInvalidInput
	codeUnsupportedPlatform
)

func errorString(code errCode, inputs ...interface{}) string {
	format, ok := errFmtStrings[currentLanguage][code]
	if !ok {
		format, _ = errFmtStrings[English][code]
	}
	return fmt.Sprintf(format, inputs...)
}

var errFmtStrings = map[Language]map[errCode]string{
	English: {
		codeNotFound:                   "%v was not found",
		codeExistingElement:            "%1v %2v already defined",
		codeExistingElementOverwritten: "%1v %2v already defined, old %2v overwritten",
		codeInsufficientInputs:         "Must supply at least %v %v",
		codeUnsupportedFormat:          "Unsupported format: %v",
		codeNilInput:                   "%v cannot be nil",
		codeIndivisibleInput:           "%v was not divisible by %v",
		codeInvalidInput:               "invalid input: %v",
		codeUnsupportedPlatform:        "%v is not supported on this platform",
	},
	Deutsch: {
		codeNotFound:                   "%v nicht gefunden",
		codeExistingElement:            "%1v %2v schon definiert",
		codeExistingElementOverwritten: "%1v %2v schon definiert, alterer %2v uberschreiben",
		codeInsufficientInputs:         "%v %v gebrauchen",
		codeUnsupportedFormat:          "Format nicht unterstützt: %v",
		codeNilInput:                   "%v darf nicht nil sein",
		codeIndivisibleInput:           "%v ist nicht teilbar durch %v",
		codeInvalidInput:               "ungültige Eingabe: %v",
		codeUnsupportedPlatform:        "%v ist auf diesem betriebssystem nicht unterstützt",
	},
	日本語: {
		codeNotFound: "%qが見つからない",
		//codeExistingElement: "%1q%2qはもう存在します"
		//codeExistingElementOverwritten: "%1q%2qはすでに存在し、古い%2qは上書きされます"
		//codeInsufficientInputs: "%v%qが要つる",
		//codeUnsupportedFormat: "対応プロトコルがない:%q",
		//codeNilInput: "%qはnilであってはなりません",
		//codeIndivisibleInput: "%qが%vに割り切れない",
		// These are commented out because I am not confident they are correct
	},
}
