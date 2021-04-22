package apperror

const (
	InvalidData            = "Os dados informados não são válidos"
	InvalidEmailAndPhone   = "Informe um email ou número de telefone válido."
	InvalidCode            = "Código inválido"
	InvalidSession         = "Sessão inválida"
	InternalError          = "Erro interno desconhecido"
	AlreadyExisting        = "Este parâmetro/recurso já está em uso"
	TooManyRequest         = "Bloqueado temporariamente por excesso de tentativa"
	PreconditionRequired   = "Pré-requisito necessário"
	Unauthorized           = "Não autorizado"
	NotFound               = "Não encontrado"
	UserNotFount           = "Usuário não encontrado"
	Forbidden              = "Operação não autorizada, atingiu o limite definido"
	ExternalError          = "Erro em recurso externo"
	MissingInternalInfo    = "Falta de informação interna"
	PermissionError        = "O usuário não tem o nível de permissão necessário"
	NoAccessPass           = "Não existe nenhuma conta vinculada a este passe de acesso."
	NotPublic              = "O recurso não está público"
	MaxFileSize            = "Limite máximo do arquivo é de 2.5 MB"
	AllParts               = "O servidor já possui todas as partes do arquivo."
	IncompleteFileOnServer = "O arquivo está incompleto no servidor."
)

type AppError interface {
	Error() string
	StatusCode() int
	Message() string
}

type ErrorTracking struct {
	Cause         error
	Status        int
	PublicMessage string
}

func (err ErrorTracking) Error() string {
	return err.Cause.Error()
}

func (err *ErrorTracking) StatusCode() int {
	return err.Status
}

func (err *ErrorTracking) Message() string {
	return err.PublicMessage
}

func NewError(cause error, status int, message string) AppError {
	return &ErrorTracking{
		Cause:         cause,
		Status:        status,
		PublicMessage: message,
	}
}
