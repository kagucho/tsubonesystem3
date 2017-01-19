package member

import (
	"net/http"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/common"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/context"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/token/authorizer"
)

func DeleteServeHTTP(writer http.ResponseWriter, request *http.Request, context context.Context, claim authorizer.Claim) {
	if context.DB.DeleteMember(request.FormValue("id")) == nil {
		common.ServeJSON(writer, struct{}{}, http.StatusOK)
	} else {
		common.ServeErrorDefault(writer, http.StatusBadRequest)
	}
}
