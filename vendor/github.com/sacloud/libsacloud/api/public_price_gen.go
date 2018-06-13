package api

/************************************************
  generated by IDE. for [PublicPriceAPI]
************************************************/

import (
	"github.com/sacloud/libsacloud/sacloud"
)

/************************************************
   To support fluent interface for Find()
************************************************/

// Reset 検索条件のリセット
func (api *PublicPriceAPI) Reset() *PublicPriceAPI {
	api.reset()
	return api
}

// Offset オフセット
func (api *PublicPriceAPI) Offset(offset int) *PublicPriceAPI {
	api.offset(offset)
	return api
}

// Limit リミット
func (api *PublicPriceAPI) Limit(limit int) *PublicPriceAPI {
	api.limit(limit)
	return api
}

// Include 取得する項目
func (api *PublicPriceAPI) Include(key string) *PublicPriceAPI {
	api.include(key)
	return api
}

// Exclude 除外する項目
func (api *PublicPriceAPI) Exclude(key string) *PublicPriceAPI {
	api.exclude(key)
	return api
}

// FilterBy 指定キーでのフィルター
func (api *PublicPriceAPI) FilterBy(key string, value interface{}) *PublicPriceAPI {
	api.filterBy(key, value, false)
	return api
}

// FilterMultiBy 任意項目でのフィルタ(完全一致 OR条件)
func (api *PublicPriceAPI) FilterMultiBy(key string, value interface{}) *PublicPriceAPI {
	api.filterBy(key, value, true)
	return api
}

// WithNameLike 名称条件(DisplayName)
func (api *PublicPriceAPI) WithNameLike(name string) *PublicPriceAPI {
	return api.FilterBy("DisplayName", name)
}

//// WithTag
//func (api *PublicPriceAPI) WithTag(tag string) *PublicPriceAPI {
//	return api.FilterBy("Tags.Name", tag)
//}

//// WithTags
//func (api *PublicPriceAPI) WithTags(tags []string) *PublicPriceAPI {
//	return api.FilterBy("Tags.Name", []interface{}{tags})
//}

// func (api *PublicPriceAPI) WithSizeGib(size int) *PublicPriceAPI {
// 	api.FilterBy("SizeMB", size*1024)
// 	return api
// }

// func (api *PublicPriceAPI) WithSharedScope() *PublicPriceAPI {
// 	api.FilterBy("Scope", "shared")
// 	return api
// }

// func (api *PublicPriceAPI) WithUserScope() *PublicPriceAPI {
// 	api.FilterBy("Scope", "user")
// 	return api
// }

// SortBy 指定キーでのソート
func (api *PublicPriceAPI) SortBy(key string, reverse bool) *PublicPriceAPI {
	api.sortBy(key, reverse)
	return api
}

// SortByName 名称でのソート(DisplayName)
func (api *PublicPriceAPI) SortByName(reverse bool) *PublicPriceAPI {
	api.sortBy("DisplayName", reverse)
	return api
}

// func (api *PublicPriceAPI) SortBySize(reverse bool) *PublicPriceAPI {
// 	api.sortBy("SizeMB", reverse)
// 	return api
// }

/************************************************
   To support Setxxx interface for Find()
************************************************/

// SetEmpty 検索条件のリセット
func (api *PublicPriceAPI) SetEmpty() {
	api.reset()
}

// SetOffset オフセット
func (api *PublicPriceAPI) SetOffset(offset int) {
	api.offset(offset)
}

// SetLimit リミット
func (api *PublicPriceAPI) SetLimit(limit int) {
	api.limit(limit)
}

// SetInclude 取得する項目
func (api *PublicPriceAPI) SetInclude(key string) {
	api.include(key)
}

// SetExclude 除外する項目
func (api *PublicPriceAPI) SetExclude(key string) {
	api.exclude(key)
}

// SetFilterBy 指定キーでのフィルター
func (api *PublicPriceAPI) SetFilterBy(key string, value interface{}) {
	api.filterBy(key, value, false)
}

// SetFilterMultiBy 任意項目でのフィルタ(完全一致 OR条件)
func (api *PublicPriceAPI) SetFilterMultiBy(key string, value interface{}) {
	api.filterBy(key, value, true)
}

// SetNameLike 名称条件(DisplayName)
func (api *PublicPriceAPI) SetNameLike(name string) {
	api.FilterBy("DisplayName", name)
}

//// SetTag
//func (api *PublicPriceAPI) SetTag(tag string) {
//}

//// SetTags
//func (api *PublicPriceAPI) SetTags(tags []string) {
//}

// func (api *PublicPriceAPI) SetSizeGib(size int) {
// 	api.FilterBy("SizeMB", size*1024)
// }

// func (api *PublicPriceAPI) SetSharedScope() {
// 	api.FilterBy("Scope", "shared")
// }

// func (api *PublicPriceAPI) SetUserScope() {
// 	api.FilterBy("Scope", "user")
// }

// SetSortBy 指定キーでのソート
func (api *PublicPriceAPI) SetSortBy(key string, reverse bool) {
	api.sortBy(key, reverse)
}

// SetSortByName 名称でのソート(DisplayName)
func (api *PublicPriceAPI) SetSortByName(reverse bool) {
	api.sortBy("DisplayName", reverse)
}

// func (api *PublicPriceAPI) SetSortBySize(reverse bool) {
// 	api.sortBy("SizeMB", reverse)
// }

/************************************************
  To support CRUD(Create/Read/Update/Delete)
************************************************/

// func (api *PublicPriceAPI) New() *sacloud.PublicPrice {
// 	return &sacloud.PublicPrice{}
// }

// func (api *PublicPriceAPI) Create(value *sacloud.PublicPrice) (*sacloud.PublicPrice, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.create(api.createRequest(value), res)
// 	})
// }

// func (api *PublicPriceAPI) Read(id int64) (*sacloud.PublicPrice, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.read(id, nil, res)
// 	})
// }

// func (api *PublicPriceAPI) Update(id int64, value *sacloud.PublicPrice) (*sacloud.PublicPrice, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.update(id, api.createRequest(value), res)
// 	})
// }

// func (api *PublicPriceAPI) Delete(id int64) (*sacloud.PublicPrice, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.delete(id, nil, res)
// 	})
// }

/************************************************
  Inner functions
************************************************/

func (api *PublicPriceAPI) setStateValue(setFunc func(*sacloud.Request)) *PublicPriceAPI {
	api.baseAPI.setStateValue(setFunc)
	return api
}

//func (api *PublicPriceAPI) request(f func(*sacloud.Response) error) (*sacloud.PublicPrice, error) {
//	res := &sacloud.Response{}
//	err := f(res)
//	if err != nil {
//		return nil, err
//	}
//	return res.ServiceClass, nil
//}
//
//func (api *PublicPriceAPI) createRequest(value *sacloud.PublicPrice) *sacloud.Request {
//	req := &sacloud.Request{}
//	req.ServiceClass = value
//	return req
//}
