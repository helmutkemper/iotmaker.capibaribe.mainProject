package main

import (
  sProxy "github.com/helmutkemper/SimpleReverseProxy"
  "net/http"
  "io/ioutil"
      "fmt"
  "github.com/helmutkemper/dockerManager/image"
  "github.com/helmutkemper/blog/config"
  tk "github.com/helmutkemper/telerik"
  "github.com/helmutkemper/dockerManager/container"
)

func containerListHtml(w sProxy.ProxyResponseWriter, r *sProxy.ProxyRequest) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")

  page, err := ioutil.ReadFile("containerStatus.html")
  if err != nil {
    w.Write( []byte( err.Error() ) )
    return
  }

  w.Write( page )
}
// fixme: quando o usuário for criar um banco de dados ou parecido, ele tem que ser avisado de exportar o diretório de dados para a máquina
func main() {
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.ProxyNotFound )
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.ProxyError )
  sProxy.FuncMap.Add( containerListHtml )
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.RouteAdd )
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.RouteDelete )
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.ProxyStatistics )
  sProxy.FuncMap.Add( image.WebList )

  sProxy.ProxyRootConfig = sProxy.ProxyConfig{
    ListenAndServe: ":8888",
  }
  err := sProxy.ProxyRootConfig.AddRouteToProxyStt(
    sProxy.ProxyRoute{
    // docker run -d --name ghost-blog-demo -p 2368:2368 ghost
      Name: "blog",
      Domain: sProxy.ProxyDomain{
        Host: "blog.localhost:8888",
      },
        ProxyEnable: true,
        ProxyServers: []sProxy.ProxyUrl{
        {
          Name: "docker 1 - ok",
          Url: "http://localhost:2368",
        },
      },
    },
  )
  err = sProxy.ProxyRootConfig.AddRouteFromFuncStt(
    sProxy.ProxyRoute{
      Name: "image raw",
      Domain: sProxy.ProxyDomain{
        NotFoundHandle: sProxy.ProxyRootConfig.ProxyNotFound,
        ErrorHandle: sProxy.ProxyRootConfig.ProxyError,
        Host: "image.localhost:8888",
      },
      Path: sProxy.ProxyPath{
        Path: "/raw",
        Method: "GET",
        ExpReg: "^/raw/(?P<id>[0-9a-fA-F]+)$",
      },
      ProxyEnable: false,
      Handle: sProxy.ProxyHandle{
        Handle: image.WebImageInfoList,
      },
    },
  )
  err = sProxy.ProxyRootConfig.AddRouteFromFuncStt(
    sProxy.ProxyRoute{
      Name: "image list",
      Domain: sProxy.ProxyDomain{
        NotFoundHandle: sProxy.ProxyRootConfig.ProxyNotFound,
        ErrorHandle: sProxy.ProxyRootConfig.ProxyError,
        Host: "",
      },
      Path: sProxy.ProxyPath{
        Path: "/image.list",
        Method: "GET",
      },
      ProxyEnable: false,
      Handle: sProxy.ProxyHandle{
        Handle: image.WebList,
      },
    },
  )
  if err != "" {
    fmt.Println( err )
  }
  err = sProxy.ProxyRootConfig.AddRouteFromFuncStt(
    sProxy.ProxyRoute{
      Name: "container list",
      Domain: sProxy.ProxyDomain{
        NotFoundHandle: sProxy.ProxyRootConfig.ProxyNotFound,
        ErrorHandle: sProxy.ProxyRootConfig.ProxyError,
        Host: "",
      },
      Path: sProxy.ProxyPath{
        Path: "/container.list",
        Method: "GET",
      },
      ProxyEnable: false,
      Handle: sProxy.ProxyHandle{
        Handle: container.WebList,
      },
    },
  )
  if err != "" {
    fmt.Println( err )
  }
  err = sProxy.ProxyRootConfig.AddRouteFromFuncStt(
    sProxy.ProxyRoute{
      Name: "image grid",
      Domain: sProxy.ProxyDomain{
        NotFoundHandle: sProxy.ProxyRootConfig.ProxyNotFound,
        ErrorHandle: sProxy.ProxyRootConfig.ProxyError,
        Host: "",
      },
      Path: sProxy.ProxyPath{
        Path: "/image.view",
        Method: "GET",
      },
      ProxyEnable: false,
      Handle: sProxy.ProxyHandle{
        Handle: func(w sProxy.ProxyResponseWriter, r *sProxy.ProxyRequest){

          data := config.Data{}
          data.LoadFromFile( "./blogData/grid.json" )

          el := tk.KendoUiGrid{
            Html: tk.HtmlElementDiv{
              Global: tk.HtmlGlobalAttributes{
                Id: "grid",
                Style: "width: 980px;", //fixme: constante no código
              },
            },
            Columns: []tk.KendoGridColumns{
              //{ Selectable: tk.TRUE, Width: 50 },
              {
                Field: "Id",
                /*Editor: tk.JavaScript{
                  Code: "function(a,b){console.log('a:', a);console.log('b:', b);}",
                },*/
              },
              { Field: "Created", Format: "{0:MM/dd/yyyy}" },
              { Field: "RepoDigests" },
              { Field: "RepoTags" },
              { Field: "Size" },
              { Field: "VirtualSize" },
              {
                Command: []tk.KendoGridColumnsCommand{
                  {
                    Name: tk.COLUMNS_COMMAND_DESTROY,
                    Text: "remove",
                  },
                  {
                    Name: tk.COLUMNS_COMMAND_EDIT,
                    Text: tk.KendoGridMessagesCommands{
                      Update: "#force_empity#",
                      Cancel: "cancel",
                    },
                    IconClass: tk.KendoGridColumnsIconClass{
                      Update: "none",
                    },
                  },
                  {
                    Name: tk.COLUMNS_COMMAND_CUSTOM,
                    Text: "view",
                  },
                },
              },
            },
            Sortable: tk.TRUE,
            PersistSelection: tk.TRUE,
            //Editable: tk.TRUE,
            Editable: tk.KendoGridEditable{
              Mode: tk.KENDO_GRID_EDITOR_MODE_POPUP,
              Confirmation: "Are you sure that you want to delete this image from server?",
            },
            ColumnMenu: tk.KendoGridColumnMenu{
              Columns: tk.FALSE,
            },
            Filterable: tk.KendoGridFilterable{
              Messages: tk.KendoGridFilterableMessages{
                And: "and",
                Or: "or",
              },
              Operators: tk.KendoGridFilterableOperators{
                String: tk.KendoGridFilterableOperatorsString{
                  IsNotEmpty: "Is not empty",
                  IsNull: "Is null",
                  IsEmpty: "Is empty",
                  DoesNotContain: "Does not contain",
                  Contains: "Contains",
                  EndsWith: "Ends with",
                  StartsWith: "Starts with",
                  Eq: "Is equal to",
                  Neq: "Is not equal to",
                },
                Date: tk.KendoGridFilterableOperatorsDate{
                  Neq: "Is not equal to",
                  IsNull: "Is null",
                  Eq: "Is equal to",
                  Gt: "Is after",
                  Lte: "Is before or equal to",
                  Lt: "Is before",
                  Gte: "Is after or equal to",
                },
                Number: tk.KendoGridFilterableOperatorsNumber{
                  Neq: "Is not equal to",
                  Lt: "Is less than",
                  Eq: "Is equal to",
                  IsNull: "Is null",
                  Lte: "Is less than or equal to",
                  Gt: "Is greater than",
                  Gte: "Is greater than or equal to",
                  IsNotNull: "Is not null",
                },
                Enums: tk.KendoGridFilterableOperatorsEnums{
                  Eq: "Is equal to",
                  IsNull: "Is null",
                  IsNotNull: "Is not null",
                  Neq: "Is not equal to",
                },
              },
            },
            Groupable: tk.KendoGridGroupable{
              ShowFooter: tk.TRUE,
              Enabled: tk.TRUE,
              Messages: tk.KendoGridGroupableMessages{
                Empty: "Drop columns here",
              },
            },
            Pageable: tk.KendoGridPageable{
              ButtonCount: 5,
              Refresh: tk.TRUE,
              Numeric: tk.TRUE,
              AlwaysVisible: tk.TRUE,
              Info: tk.TRUE,
              Messages: tk.KendoGridPageableMessages{
                Display: "Showing {0}-{1} from {2} data items",
                Empty: "No data",
                Page: "Enter page",
                Of: "from {0}",
                ItemsPerPage: "data items per page",
                First: "First page",
                Last: "Last page",
                Next: "Next page",
                Previous: "Previous page",
                Refresh: "Refresh the grid",
                MorePages: "More pages",
              },
            },
            DataSource: tk.KendoDataSource{
              //Type: KENDO_TYPE_DATA_JSON,
              Transport: tk.KendoTransport{
                Read: tk.KendoRead{
                  Url: "http://localhost:8888/image.list",
                  Type: tk.HTML_METHOD_GET,
                  DataType: tk.KENDO_TYPE_DATA_JSON_JSON,
                },
              },
              Page: 1,
              PageSize: 10,
              Schema: tk.KendoSchema{
                Data:  "Objects",
                Total: "Meta.TotalCount",
                Parser: tk.JavaScript{
                  Code: `function(data){
                    for( var i = 0, l = data.Objects.length; i < l; i += 1 ){
                      data.Objects[i].Created = new Date(data.Objects[i].Created * 1000);
                      if( Array.isArray( data.Objects[i].RepoDigests ) == true ){
                        var out = "";
                        for(var k = 0, lk = data.Objects[i].RepoDigests.length; k < lk; k += 1){
                          if(k != 0){
                            out += "; ";
                          }
                          out += data.Objects[i].RepoDigests[k]
                        }
                        data.Objects[i].RepoDigests = out;
                      }
                      if( Array.isArray( data.Objects[i].RepoTags ) == true ){
                        var out = "";
                        for(var k = 0, lk = data.Objects[i].RepoTags.length; k < lk; k += 1){
                          if(k != 0){
                            out += "; ";
                          }
                          out += data.Objects[i].RepoTags[k]
                        }
                        data.Objects[i].RepoTags = out;
                      }
                      if( data.Objects[i].Size >= 1024 * 1024 * 1024 ){
                        data.Objects[i].Size /= 1024 * 1024 * 1024;
                        data.Objects[i].Size = Number.parseFloat( data.Objects[i].Size ).toPrecision(4);
                        data.Objects[i].Size = '' + data.Objects[i].Size + 'GB';
                      } else if( data.Objects[i].Size >= 1024 * 1024 ){
                        data.Objects[i].Size /= 1024 * 1024;
                        data.Objects[i].Size = Number.parseFloat( data.Objects[i].Size ).toPrecision(4);
                        data.Objects[i].Size = '' + data.Objects[i].Size + 'MB';
                      } else if( data.Objects[i].Size >= 1024 ){
                        data.Objects[i].Size /= 1024;
                        data.Objects[i].Size = Number.parseFloat( data.Objects[i].Size ).toPrecision(4);
                        data.Objects[i].Size = '' + data.Objects[i].Size + 'KB';
                      } else {
                        data.Objects[i].Size = '' + data.Objects[i].Size + 'B';
                      }
                      
                      if( data.Objects[i].VirtualSize >= 1024 * 1024 * 1024 ){
                        data.Objects[i].VirtualSize /= 1024 * 1024 * 1024;
                        data.Objects[i].VirtualSize = Number.parseFloat( data.Objects[i].VirtualSize ).toPrecision(4);
                        data.Objects[i].VirtualSize += 'GB';
                      } else if( data.Objects[i].VirtualSize >= 1024 * 1024 ){
                        data.Objects[i].VirtualSize /= 1024 * 1024;
                        data.Objects[i].VirtualSize = Number.parseFloat( data.Objects[i].VirtualSize ).toPrecision(4);
                        data.Objects[i].VirtualSize += 'MB';
                      } else if( data.Objects[i].VirtualSize >= 1024 ){
                        data.Objects[i].VirtualSize /= 1024;
                        data.Objects[i].VirtualSize = Number.parseFloat( data.Objects[i].VirtualSize ).toPrecision(4);
                        data.Objects[i].VirtualSize += 'KB';
                      } else {
                        data.Objects[i].VirtualSize += 'B';
                      }
                    }
                    return data;
                  }`,
                },
                Model: tk.KendoDataModel{
                  Id: "Id",
                  Fields: map[string]tk.KendoField{
                    "Id": {
                      Type: tk.JAVASCRIPT_STRING,
                    },
                    "Created": {
                      Type: tk.JAVASCRIPT_DATE,
                    },
                    "RepoTags": {
                      Type: tk.JAVASCRIPT_STRING,
                    },
                    "RepoDigests": {
                      Type: tk.JAVASCRIPT_STRING,
                    },
                    "Size": {
                      Type: tk.JAVASCRIPT_STRING,
                    },
                    "VirtualSize": {
                      Type: tk.JAVASCRIPT_STRING,
                    },
                  },
                },
              },
            },
          }

          data.TelerikOnLoadCode = string( el.ToJavaScript() )
          data.Post[0].Text = string( el.ToHtml() )

          data.TemplateToServer("./static/template", w)

        },
      },
    },
  )
  if err != "" {
    fmt.Println( err )
  }

  sProxy.ProxyRootConfig.Prepare()

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer( http.Dir( "static" ) ) ) )
  http.HandleFunc("/", sProxy.ProxyFunc)
  http.ListenAndServe(sProxy.ProxyRootConfig.ListenAndServe, nil)
}