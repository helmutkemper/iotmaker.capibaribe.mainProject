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

          el := tk.HtmlElementDiv{
            Global: tk.HtmlGlobalAttributes{
              Id: "main",
              Class: "k-content",
            },
            Content: tk.Content{

              &tk.KendoUiGrid{
                Html: tk.HtmlElementDiv{
                  Global: tk.HtmlGlobalAttributes{
                    Id: "grid",
                  },
                },
                Columns: []tk.KendoGridColumns{
                  //{ Selectable: tk.TRUE, Width: 50 },
                  {
                    Field: "Id",
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
                        Name: tk.COLUMNS_COMMAND_VIEW,
                        Text: "view",
                        IconClass: tk.ICON_PREVIEW,

                        ViewWindow: &tk.KendoUiWindow{
                          Html: tk.HtmlElementDiv{
                            Global: tk.HtmlGlobalAttributes{
                              Id: "imageDetailsWindow",
                            },
                          },
                          Title: "Image Data",
                          Modal: tk.TRUE,
                          Visible: tk.FALSE,
                          Resizable: tk.TRUE,
                          Width: 900,

                        },

                        ViewTemplate: &tk.HtmlElementScript{
                          Global: tk.HtmlGlobalAttributes{
                            Id: "template-details-container",
                          },
                          Type: tk.SCRIPT_TYPE_KENDO_TEMPLATE,
                          Content: tk.Content{

                            &tk.HtmlElementDiv{
                              Global: tk.HtmlGlobalAttributes{
                                Id: "details-container",
                              },
                              Content: tk.Content{

                                `<h3>#= Id #</h3>
                                <dl>
                                  <dt>Created:</dt><dt>#= kendo.toString(Created, "MM/dd/yyyy") #</dt>
                                  <dt>Tags:</dt><dt>#= RepoTags #</dt>
                                  <dt>Digests:</dt><dt>#= RepoDigests #</dt>
                                  <dt>Size:</dt><dt>#= Size #</dt>
                                  <dt>Virtual Size:</dt><dt>#= VirtualSize #</dt>
                                </dl>`,
                              },
                            },

                          },
                        },

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
              },

            },

          }

          data.TelerikVarGlobal = string( el.Content.MakeJsObject() )
          data.TelerikOnLoadCode = string( el.Content.ToJavaScript() )

          data.TelerikScriptTemplate = string( el.Content.MakeJsScript() )
          data.TelerikHtmlSupport = string( el.ToHtmlSupport() )

          data.Post[0].Title = "Image List"
          data.Post[0].Text += string( el.ToHtml() )

          data.TemplateToServer("./static/template", w)

        },
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
        Path: "/image.edit",
        Method: "GET",
        ExpReg: "^/image.edit/(?P<id>[0-9a-fA-F]+)$",
      },
      ProxyEnable: false,
      Handle: sProxy.ProxyHandle{
        Handle: func(w sProxy.ProxyResponseWriter, r *sProxy.ProxyRequest){

          data := config.Data{}
          data.LoadFromFile( "./blogData/grid.json" )

          el := tk.HtmlElementDiv{
            Global: tk.HtmlGlobalAttributes{
              Id: "spanCreateTemplateExposedPortsAddNewPort",
              Class: "k-content",
              Style: "width: 300px !important;",
            },
            Content: tk.Content{

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigHostName",
                    Content: tk.Content{
                      "Host Name",
                    },
                    Global: tk.HtmlGlobalAttributes{
                      Style: "width: 200px important!;",
                    },
                  },

                  &tk.HtmlInputText{
                    Global: tk.HtmlGlobalAttributes{
                      Id: "ConfigHostName",
                      Class: "k-textbox",
                    },
                    Name: "HostName",
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigDomainName",
                    Content: tk.Content{
                      "Domain Name",
                    },
                  },

                  &tk.HtmlInputText{
                    Global: tk.HtmlGlobalAttributes{
                      Id: "ConfigDomainName",
                      Class: "k-textbox",
                    },
                    Name: "DomainName",
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigUser",
                    Content: tk.Content{
                      "User",
                    },
                  },

                  &tk.HtmlInputText{
                    Global: tk.HtmlGlobalAttributes{
                      Id: "ConfigUser",
                      Class: "k-textbox",
                    },
                    Name: "User",
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigAttachStdIn",
                    Content: tk.Content{
                      "Attach Std In",
                    },
                  },

                  &tk.KendoUiDropDownList{
                    Html: tk.HtmlInputText{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigAttachStdIn",
                      },
                      Name: "AttachStdIn",
                    },
                    DataValueField: "key",
                    DataTextField: "value",

                    DataSource: []map[string]interface{}{
                      {
                        "key": -1,
                        "value": "Default",
                      },
                      {
                        "key": 0,
                        "value": "No",
                      },
                      {
                        "key": 1,
                        "value": "Yes",
                      },
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigAttachStdOut",
                    Content: tk.Content{
                      "Attach Std Out",
                    },
                  },

                  &tk.KendoUiDropDownList{
                    Html: tk.HtmlInputText{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigAttachStdOut",
                      },
                      Name: "AttachStdOut",
                    },
                    DataValueField: "key",
                    DataTextField: "value",

                    DataSource: []map[string]interface{}{
                      {
                        "key": -1,
                        "value": "Default",
                      },
                      {
                        "key": 0,
                        "value": "No",
                      },
                      {
                        "key": 1,
                        "value": "Yes",
                      },
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigAttachStdErr",
                    Content: tk.Content{
                      "Host Attach Std Err",
                    },
                  },

                  &tk.KendoUiDropDownList{
                    Html: tk.HtmlInputText{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigAttachStdErr",
                      },
                      Name: "AttachStdErr",
                    },
                    DataValueField: "key",
                    DataTextField: "value",

                    DataSource: []map[string]interface{}{
                      {
                        "key": -1,
                        "value": "Default",
                      },
                      {
                        "key": 0,
                        "value": "No",
                      },
                      {
                        "key": 1,
                        "value": "Yes",
                      },
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigExposedPorts",
                    Content: tk.Content{
                      "Exposed Ports",
                    },
                  },

                  &tk.KendoUiMultiSelect{
                    Html: tk.HtmlElementFormSelect{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigExposedPorts",
                      },
                      Name: "ExposedPorts",
                    },
                    ClearButton: tk.FALSE,
                    DataValueField: "id",
                    DataTextField: "ExposedPortsShow",
                    DataSource: tk.KendoDataSource{
                      //Type: KENDO_TYPE_DATA_JSON,
                      Transport: tk.KendoTransport{
                        Read: tk.KendoRead{
                          Url: "/static/test/read",
                          Type: tk.HTML_METHOD_GET,
                          DataType: tk.KENDO_TYPE_DATA_JSON_JSON,
                        },
                        Create: tk.KendoCreate{
                          Url: "/static/test/create",
                          Type: tk.HTML_METHOD_POST,
                          DataType: tk.KENDO_TYPE_DATA_JSON_JSON,
                        },
                      },
                      Schema: tk.KendoSchema{
                        Data:  "Objects",
                        Total: "Total",
                        Model: tk.KendoDataModel{
                          Id: "id",
                          Fields: map[string]tk.KendoField{
                            "id": {
                              Type: tk.JAVASCRIPT_NUMBER,
                            },
                            "ExposedPortsNumber": {
                              Type: tk.JAVASCRIPT_NUMBER,
                            },
                            "ExposedPortsProtocol": {
                              Type: tk.JAVASCRIPT_STRING,
                            },
                            "ExposedPortsShow": {
                              Type: tk.JAVASCRIPT_STRING,
                            },
                          },
                        },
                      },
                      //PageSize: 10,
                      ServerPaging: tk.TRUE,
                    },
                    Dialog: tk.KendoUiDialog{
                      Html: tk.HtmlElementDiv{
                        Global: tk.HtmlGlobalAttributes{
                          Id: tk.GetAutoId(),
                        },
                      },
                      Title: "Add new exposed port.",
                      Content: tk.Content{

                        // regra, o form valida automaticamente
                        &tk.HtmlElementDiv{
                          Global: tk.HtmlGlobalAttributes{
                            Id: "ConfigExposedPortsDialogContent",
                            Class: "k-content",
                          },
                          Content: tk.Content{

                            &tk.HtmlElementDiv{
                              Content: tk.Content{

                                &tk.HtmlElementFormLabel{
                                  For: "ExposedPortsNumber",
                                  Content: tk.Content{
                                    "Port number",
                                  },
                                },

                                &tk.KendoUiNumericTextBox{
                                  Html: tk.HtmlInputNumber{
                                    Name:         "ExposedPortsNumber",
                                    PlaceHolder:  "",
                                    AutoComplete: tk.FALSE,
                                    Required:     tk.TRUE,
                                    // Pattern: "[^=]*",
                                    Global: tk.HtmlGlobalAttributes{
                                      Id:    "ExposedPortsNumber",
                                      Class: "oneThirdSize",
                                      Extra: map[string]interface{}{
                                        "validationMessage": "Enter a {0}",
                                      },
                                    },
                                  },
                                  Format: "#",
                                },
                              },
                            },

                            &tk.HtmlElementDiv{
                              Content: tk.Content{

                                &tk.HtmlElementFormLabel{
                                  For: "ExposedPortsProtocol",
                                  Content: tk.Content{
                                    "Port protocol",
                                  },
                                },

                                &tk.KendoUiComboBox{
                                  Html: tk.HtmlElementFormSelect{
                                    Global: tk.HtmlGlobalAttributes{
                                      Id:    "ExposedPortsProtocol",
                                      Class: "oneThirdSize",
                                      Data:  map[string]string{"required-msg": "Select start time"},
                                    },
                                    Name:     "ExposedPortsProtocol",
                                    Required: tk.TRUE,
                                    Options: []tk.HtmlOptions{
                                      {
                                        Label: "Please, select one",
                                        Key:   "",
                                      },
                                      {
                                        Label: "TCP",
                                        Key:   "TCP",
                                      },
                                      {
                                        Label: "UDP",
                                        Key:   "UDP",
                                      },
                                    },
                                  },
                                },
                              },
                            },

                            &tk.HtmlInputHidden{
                              Global: tk.HtmlGlobalAttributes{
                                Id: "ExposedPortsShow",
                              },
                              Name: "ExposedPortsShow",
                            },

                          },
                        },
                      },
                      Visible: tk.FALSE,
                      Width: 400,
                      Actions: []tk.KendoActions{
                        {
                          Primary: tk.FALSE,
                          Text:    "Close",
                        },
                        {
                          Primary: tk.FALSE,
                          Text:    "Add",
                          ButtonType: tk.BUTTON_TYPE_ADD,
                        },
                        {
                          Primary: tk.TRUE,
                          Text:    "Add and close",
                          ButtonType: tk.BUTTON_TYPE_ADD_AND_CLOSE,
                        },
                      },
                    },
                    NoDataTemplate: tk.HtmlElementScript{
                      Global: tk.HtmlGlobalAttributes{
                        Id: tk.GetAutoId(),
                      },
                      Type: tk.SCRIPT_TYPE_KENDO_TEMPLATE,
                      Content: tk.Content{

                        &tk.HtmlElementDiv{
                          Content: tk.Content{
                            "No data found. Do you want to add new item?",
                          },
                        },

                        "<br>",
                        "<br>",

                        &tk.HtmlElementFormButton{
                          ButtonType: tk.BUTTON_TYPE_ADD_IN_TEMPLATE,
                          Global: tk.HtmlGlobalAttributes{
                            Id: tk.GetAutoId(),
                            Class: "k-button",
                          },
                          Content: tk.Content{
                            "Add new item",
                          },
                        },

                      },
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigTry",
                    Content: tk.Content{
                      "Try",
                    },
                  },

                  &tk.KendoUiDropDownList{
                    Html: tk.HtmlInputText{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigTry",
                      },
                      Name: "Try",
                    },
                    DataValueField: "key",
                    DataTextField: "value",

                    DataSource: []map[string]interface{}{
                      {
                        "key": -1,
                        "value": "Default",
                      },
                      {
                        "key": 0,
                        "value": "No",
                      },
                      {
                        "key": 1,
                        "value": "Yes",
                      },
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigOpenStdIn",
                    Content: tk.Content{
                      "Open Std In",
                    },
                  },

                  &tk.KendoUiDropDownList{
                    Html: tk.HtmlInputText{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigOpenStdIn",
                      },
                      Name: "OpenStdIn",
                    },
                    DataValueField: "key",
                    DataTextField: "value",

                    DataSource: []map[string]interface{}{
                      {
                        "key": -1,
                        "value": "Default",
                      },
                      {
                        "key": 0,
                        "value": "No",
                      },
                      {
                        "key": 1,
                        "value": "Yes",
                      },
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigStdInOnce",
                    Content: tk.Content{
                      "Std In Once",
                    },
                  },

                  &tk.KendoUiDropDownList{
                    Html: tk.HtmlInputText{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigStdInOnce",
                      },
                      Name: "StdInOnce",
                    },
                    DataValueField: "key",
                    DataTextField: "value",

                    DataSource: []map[string]interface{}{
                      {
                        "key": -1,
                        "value": "Default",
                      },
                      {
                        "key": 0,
                        "value": "No",
                      },
                      {
                        "key": 1,
                        "value": "Yes",
                      },
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigEnv",
                    Content: tk.Content{
                      "Env",
                    },
                  },

                  &tk.KendoUiMultiSelect{
                    Html: tk.HtmlElementFormSelect{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigEnv",
                      },
                      Name: "Env",
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigCmd",
                    Content: tk.Content{
                      "Cmd",
                    },
                  },

                  &tk.KendoUiMultiSelect{
                    Html: tk.HtmlElementFormSelect{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigCmd",
                      },
                      Name: "Cmd",
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigHealthCheck",
                    Content: tk.Content{
                      "Health Check",
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigArgsEscaped",
                    Content: tk.Content{
                      "Args Escaped",
                    },
                  },

                  &tk.KendoUiDropDownList{
                    Html: tk.HtmlInputText{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigArgsEscaped",
                      },
                      Name: "ArgsEscaped",
                    },
                    DataValueField: "key",
                    DataTextField: "value",

                    DataSource: []map[string]interface{}{
                      {
                        "key": -1,
                        "value": "Default",
                      },
                      {
                        "key": 0,
                        "value": "No",
                      },
                      {
                        "key": 1,
                        "value": "Yes",
                      },
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigImage",
                    Content: tk.Content{
                      "Image",
                    },
                  },

                  &tk.KendoUiAutoComplete{
                    Html: tk.HtmlInputText{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigHostImage",
                        Class: "k-textbox",
                      },
                      Name: "Image",
                    },
                    DataTextField: "ExposedPortsShow",
                    ClearButton: tk.FALSE,
                    DataSource: tk.KendoDataSource{
                      //VarName: "testDataSource",
                      //Type: KENDO_TYPE_DATA_JSON,
                      Transport: tk.KendoTransport{
                        Read: tk.KendoRead{
                          Url: "/static/test/read",
                          Type: tk.HTML_METHOD_GET,
                          DataType: tk.KENDO_TYPE_DATA_JSON_JSON,
                        },
                        Create: tk.KendoCreate{
                          Url: "/static/test/create",
                          Type: tk.HTML_METHOD_POST,
                          DataType: tk.KENDO_TYPE_DATA_JSON_JSON,
                        },
                      },
                      Schema: tk.KendoSchema{
                        Data:  "Objects",
                        Total: "Total",
                        Model: tk.KendoDataModel{
                          Id: "id",
                          Fields: map[string]tk.KendoField{
                            "id": {
                              Type: tk.JAVASCRIPT_NUMBER,
                            },
                            "ExposedPortsNumber": {
                              Type: tk.JAVASCRIPT_NUMBER,
                            },
                            "ExposedPortsProtocol": {
                              Type: tk.JAVASCRIPT_STRING,
                            },
                            "ExposedPortsShow": {
                              Type: tk.JAVASCRIPT_STRING,
                            },
                          },
                        },
                      },
                      //PageSize: 10,
                      ServerPaging: tk.TRUE,
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigVolumes",
                    Content: tk.Content{
                      "Volumes",
                    },
                  },

                  &tk.KendoUiMultiSelect{
                    Html: tk.HtmlElementFormSelect{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigVolumes",
                      },
                      Name: "Volumes",
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigWorkingDir",
                    Content: tk.Content{
                      "Working Dir",
                    },
                  },

                  &tk.HtmlInputText{
                    Global: tk.HtmlGlobalAttributes{
                      Id: "ConfigWorkingDir",
                      Class: "k-textbox",
                    },
                    Name: "WorkingDir",
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigEntryPoint",
                    Content: tk.Content{
                      "EntryPoint",
                    },
                  },

                  &tk.KendoUiMultiSelect{
                    Html: tk.HtmlElementFormSelect{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigEntryPoint",
                      },
                      Name: "EntryPoint",
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigNetworkDisabled",
                    Content: tk.Content{
                      "Network Disabled",
                    },
                  },

                  &tk.KendoUiDropDownList{
                    Html: tk.HtmlInputText{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigNetworkDisabled",
                      },
                      Name: "NetworkDisabled",
                    },
                    DataValueField: "key",
                    DataTextField: "value",

                    DataSource: []map[string]interface{}{
                      {
                        "key": -1,
                        "value": "Default",
                      },
                      {
                        "key": 0,
                        "value": "No",
                      },
                      {
                        "key": 1,
                        "value": "Yes",
                      },
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigMacAddress",
                    Content: tk.Content{
                      "Mac Address",
                    },
                  },

                  &tk.HtmlInputText{
                    Global: tk.HtmlGlobalAttributes{
                      Id: "ConfigMacAddress",
                      Class: "k-textbox",
                    },
                    Name: "MacAddress",
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigOnBuild",
                    Content: tk.Content{
                      "On Build",
                    },
                  },

                  &tk.KendoUiMultiSelect{
                    Html: tk.HtmlElementFormSelect{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigOnBuild",
                      },
                      Name: "OnBuild",
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigLabels",
                    Content: tk.Content{
                      "Labels",
                    },
                  },

                  &tk.KendoUiMultiSelect{
                    Html: tk.HtmlElementFormSelect{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigLabels",
                      },
                      Name: "Labels",
                    },
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigStopSignal",
                    Content: tk.Content{
                      "Stop Signal",
                    },
                  },

                  &tk.HtmlInputText{
                    Global: tk.HtmlGlobalAttributes{
                      Id: "ConfigStopSignal",
                      Class: "k-textbox",
                    },
                    Name: "StopSignal",
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigStopTimeout",
                    Content: tk.Content{
                      "Stop Timeout",
                    },
                  },

                  &tk.KendoUiNumericTextBox{
                    Html: tk.HtmlInputNumber{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigStopTimeout",
                      },
                      Name: "StopTimeout",
                    },
                    Format: "#s",
                  },

                },
              },

              &tk.HtmlElementDiv{
                Content: tk.Content{

                  &tk.HtmlElementFormLabel{
                    For: "ConfigShell",
                    Content: tk.Content{
                      "Shell",
                    },
                  },

                  &tk.KendoUiMultiSelect{
                    Html: tk.HtmlElementFormSelect{
                      Global: tk.HtmlGlobalAttributes{
                        Id: "ConfigShell",
                      },
                      Name: "Shell",
                    },
                  },

                },
              },

            },
          }

          data.TelerikVarGlobal = string( el.Content.MakeJsObject() )
          data.TelerikOnLoadCode = string( el.Content.ToJavaScript() )

          data.TelerikScriptTemplate = string( el.Content.MakeJsScript() )
          data.TelerikHtmlSupport = string( el.ToHtmlSupport() )

          data.Post[0].Title = "Image List"
          data.Post[0].Text += string( el.ToHtml() )

          data.TemplateToServer("./static/template", w)

        },
      },
    },
  )
  if err != "" {
    fmt.Println( err )
  }































































  err = sProxy.ProxyRootConfig.AddRouteFromFuncStt(
    sProxy.ProxyRoute{
      Name: "image list",
      Domain: sProxy.ProxyDomain{
        NotFoundHandle: sProxy.ProxyRootConfig.ProxyNotFound,
        ErrorHandle: sProxy.ProxyRootConfig.ProxyError,
        Host: "",
      },
      Path: sProxy.ProxyPath{
        Path: "/image.info",
        Method: "GET",
        ExpReg: "^/image.info/(?P<id>[0-9a-fA-F]+)$",
      },
      ProxyEnable: false,
      Handle: sProxy.ProxyHandle{
        Handle: image.WebImageInfoList,
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