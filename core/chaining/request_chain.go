package requestchain

import "startpoint/core/model"

func ResolveRequestChain(r *model.RequestMold, all []*model.RequestMold) []*model.RequestMold {
	chain := []*model.RequestMold{}
	chain = append(chain, resolvePreviousReq(r.PreviousReq(), all)...)
	chain = append(chain, r)
	return chain
}

func resolvePreviousReq(prevReqName string, all []*model.RequestMold) []*model.RequestMold {
	if len(prevReqName) == 0 {
		return []*model.RequestMold{}
	}

	var prevReq *model.RequestMold
	for _, v := range all {
		if v.Name == prevReqName {
			prevReq = v
			break
		}
	}

	prevReqs := []*model.RequestMold{}
	if prevReq != nil {
		ascendants := resolvePreviousReq(prevReq.PreviousReq(), all)
		prevReqs = append(prevReqs, ascendants...)
		prevReqs = append(prevReqs, prevReq)
	}

	return prevReqs
}
