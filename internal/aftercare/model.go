package aftercare

const DescriptionLimit = 100

type State string

const (
	StateLoading State = "loading"
	StateError   State = "error"
	StateResult  State = "result"
)

type Ticket struct {
	ID             string `json:"id"`
	Description    string `json:"description"`
	CharacterCount int    `json:"character_count"`
}

type Problem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type View struct {
	State   State    `json:"state"`
	Loading bool     `json:"loading"`
	Error   *Problem `json:"error,omitempty"`
	Result  *Ticket  `json:"result,omitempty"`
}

func LoadingView() View {
	return View{State: StateLoading, Loading: true}
}

func ErrorView(code, message string) View {
	return View{
		State: StateError,
		Error: &Problem{Code: code, Message: message},
	}
}

func ResultView(ticket Ticket) View {
	return View{State: StateResult, Result: &ticket}
}
