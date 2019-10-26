package es

const mapping = `{
   "settings":{
      "index.number_of_shards": 1,
      "index.number_of_replicas": 5,
      "index.mapping.nested_fields.limit":50,
      "index.requests.cache.enable":true,
      "index.mapper.dynamic":false
   },
   "mappings":{
      "_default_":{
         "_all":{
            "enabled":false
         }
      },
      "%s":%s
   }
}`

// //
// 		{
// 		 "properties":{
// 		    "fieldname":{
// 		       "type":"string",
// 		    },
// 	    }
