package main

import (
  mktp "github.com/helmutkemper/SimpleReverseProxy"
  "net/http"
  "io/ioutil"
  "html/template"
  "reflect"
)

type Icon struct {
  Rel                 string
  HRef                string
}

type Image struct {
  Alt                 string
  Src                 string
}

type Link struct {
  HRef                string
  Text                string
  Alt                 string
  Active              bool
}

type Content struct {
  Id                  string
  Title               string
  AuthorName          string
  AuthorLink          string
  Date                string
  ImageCaption        string
  Image               Image
  Description         string
  Text                string
}

type SidePanel struct {
  Id                  string
  Content             []SidePanelContent
  BreakLine           string
}

type SidePanelContent struct {
  Title               string
  Content             string
}

type PageData struct {
  BasePath            string
  Lang                string
  PageTitle           string
  AuthorInformation   string
  Favico              Icon
  Template            Template
  Content             []Content
  Copyright           string
  EditEnable          bool
  SaveButtonText      string
}

type Template struct {
  NavigationMenu      []Link
  ListGroup           []Link
  SidePanel           []SidePanel
  AutoCompleteEnable  bool
}

func blogNovo(w mktp.ProxyResponseWriter, r *mktp.ProxyRequest) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")

  tmpl, err := template.New("page_right_side.tmpl").Funcs(template.FuncMap{
    "htmlSafe": func(html string) template.HTML {
      return template.HTML(html)
    },
    "last": func(a interface{}) int {
      return reflect.ValueOf(a).Len() - 1
    },
  }).ParseFiles("./static/template/blog/page_right_side.tmpl")
  if err != nil {
    w.Write( []byte( err.Error() ) )
  }

  pageData := PageData{
    BasePath: "./static/template/blog/",
    Lang: "en",
    PageTitle: "Kemper.com.br",
    AuthorInformation: "Helmut Kemper",
    Copyright: "Helmut Kemper 2018",
    Favico: Icon{
      Rel: "shortcut icon",
      HRef: "favicon.ico",
    },
    EditEnable: false,
    SaveButtonText: "Save this article",
    Template: Template{
      AutoCompleteEnable: true,
      NavigationMenu: []Link{
        {
          HRef: "#",
          Alt: "",
          Text: "Nav Menu 1",
        },
        {
          HRef: "#",
          Alt: "",
          Text: "Nav Menu 2",
        },
        {
          HRef: "#",
          Alt: "",
          Text: "Nav Menu 3",
        },
      },
      ListGroup: []Link{
        {
          HRef: "#",
          Alt: "",
          Text: "Category 1",
          Active: true,
        },
        {
          HRef: "#",
          Alt: "",
          Text: "Category 2",
        },
        {
          HRef: "#",
          Alt: "",
          Text: "Category 3",
        },
      },
      SidePanel: []SidePanel{
        {
          Id: "SideBarPanel1",
          BreakLine: "<br><br>",
          Content:[]SidePanelContent{
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
          },
        },
        {
          Id: "SideBarPanel2",
          BreakLine: "<br><br>",
          Content:[]SidePanelContent{
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
          },
        },
        {
          Id: "SideBarPanel3",
          BreakLine: "<br><br>",
          Content:[]SidePanelContent{
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
          },
        },
        {
          Id: "SideBarPanel4",
          BreakLine: "<br><br>",
          Content:[]SidePanelContent{
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
          },
        },
        {
          Id: "SideBarPanel5",
          BreakLine: "<br><br>",
          Content:[]SidePanelContent{
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
            {
              Title: "Sidebar panel widget",
              Content: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
            },
          },
        },
      },
    },
    Content: []Content{
      {
        Id: "Post1",
        Title: "Primeiro post",
        AuthorName: "Helmut Kemper",
        AuthorLink: "#",
        Date: "12 January 2015 10:00 am",
        ImageCaption: "Caption here",
        Image: Image{
          Src: "http://placehold.it/900x400",
          Alt: "image from post",
        },
        Description: "Mussum Ipsum, cacilds vidis litro abertis. Quem manda na minha terra sou euzis! Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus. Interessantiss quisso pudia ce receita de bolis, mais bolis eu num gostis. Aenean aliquam molestie leo, vitae iaculis nisl.",
        Text: "<p>Mussum Ipsum, cacilds vidis litro abertis. Quem manda na minha terra sou euzis! Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus. Interessantiss quisso pudia ce receita de bolis, mais bolis eu num gostis. Aenean aliquam molestie leo, vitae iaculis nisl.</p><p>Nullam volutpat risus nec leo commodo, ut interdum diam laoreet. Sed non consequat odio. Per aumento de cachacis, eu reclamis. Em pé sem cair, deitado sem dormir, sentado sem cochilar e fazendo pose. Delegadis gente finis, bibendum egestas augue arcu ut est.</p><p>Mussum Ipsum, cacilds vidis litro abertis. Quem manda na minha terra sou euzis! Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus. Interessantiss quisso pudia ce receita de bolis, mais bolis eu num gostis. Aenean aliquam molestie leo, vitae iaculis nisl.</p><p>Nullam volutpat risus nec leo commodo, ut interdum diam laoreet. Sed non consequat odio. Per aumento de cachacis, eu reclamis. Em pé sem cair, deitado sem dormir, sentado sem cochilar e fazendo pose. Delegadis gente finis, bibendum egestas augue arcu ut est.</p>",
      },
      {
        Id: "Post2",
        Title: "Primeiro post",
        AuthorName: "Helmut Kemper",
        AuthorLink: "#",
        Date: "12 January 2015 10:00 am",
        ImageCaption: "Caption here",
        Image: Image{
          Src: "http://placehold.it/900x400",
          Alt: "image from post",
        },
        Description: "Mussum Ipsum, cacilds vidis litro abertis. Quem manda na minha terra sou euzis! Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus. Interessantiss quisso pudia ce receita de bolis, mais bolis eu num gostis. Aenean aliquam molestie leo, vitae iaculis nisl.",
        Text: "<p>Mussum Ipsum, cacilds vidis litro abertis. Quem manda na minha terra sou euzis! Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus. Interessantiss quisso pudia ce receita de bolis, mais bolis eu num gostis. Aenean aliquam molestie leo, vitae iaculis nisl.</p><p>Nullam volutpat risus nec leo commodo, ut interdum diam laoreet. Sed non consequat odio. Per aumento de cachacis, eu reclamis. Em pé sem cair, deitado sem dormir, sentado sem cochilar e fazendo pose. Delegadis gente finis, bibendum egestas augue arcu ut est.</p><p>Mussum Ipsum, cacilds vidis litro abertis. Quem manda na minha terra sou euzis! Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus. Interessantiss quisso pudia ce receita de bolis, mais bolis eu num gostis. Aenean aliquam molestie leo, vitae iaculis nisl.</p><p>Nullam volutpat risus nec leo commodo, ut interdum diam laoreet. Sed non consequat odio. Per aumento de cachacis, eu reclamis. Em pé sem cair, deitado sem dormir, sentado sem cochilar e fazendo pose. Delegadis gente finis, bibendum egestas augue arcu ut est.</p>",
      },
      {
        Id: "Post2",
        Title: "Primeiro post",
        AuthorName: "Helmut Kemper",
        AuthorLink: "#",
        Date: "12 January 2015 10:00 am",
        ImageCaption: "Caption here",
        Image: Image{
          Src: "http://placehold.it/900x400",
          Alt: "image from post",
        },
        Description: "Mussum Ipsum, cacilds vidis litro abertis. Quem manda na minha terra sou euzis! Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus. Interessantiss quisso pudia ce receita de bolis, mais bolis eu num gostis. Aenean aliquam molestie leo, vitae iaculis nisl.",
        Text: "<p>Mussum Ipsum, cacilds vidis litro abertis. Quem manda na minha terra sou euzis! Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus. Interessantiss quisso pudia ce receita de bolis, mais bolis eu num gostis. Aenean aliquam molestie leo, vitae iaculis nisl.</p><p>Nullam volutpat risus nec leo commodo, ut interdum diam laoreet. Sed non consequat odio. Per aumento de cachacis, eu reclamis. Em pé sem cair, deitado sem dormir, sentado sem cochilar e fazendo pose. Delegadis gente finis, bibendum egestas augue arcu ut est.</p><p>Mussum Ipsum, cacilds vidis litro abertis. Quem manda na minha terra sou euzis! Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus. Interessantiss quisso pudia ce receita de bolis, mais bolis eu num gostis. Aenean aliquam molestie leo, vitae iaculis nisl.</p><p>Nullam volutpat risus nec leo commodo, ut interdum diam laoreet. Sed non consequat odio. Per aumento de cachacis, eu reclamis. Em pé sem cair, deitado sem dormir, sentado sem cochilar e fazendo pose. Delegadis gente finis, bibendum egestas augue arcu ut est.</p>",
      },
      {
        Id: "Post4",
        Title: "Primeiro post",
        AuthorName: "Helmut Kemper",
        AuthorLink: "#",
        Date: "12 January 2015 10:00 am",
        ImageCaption: "Caption here",
        Image: Image{
          Src: "http://placehold.it/900x400",
          Alt: "image from post",
        },
        Description: "Mussum Ipsum, cacilds vidis litro abertis. Quem manda na minha terra sou euzis! Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus. Interessantiss quisso pudia ce receita de bolis, mais bolis eu num gostis. Aenean aliquam molestie leo, vitae iaculis nisl.",
        Text: "<p>Mussum Ipsum, cacilds vidis litro abertis. Quem manda na minha terra sou euzis! Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus. Interessantiss quisso pudia ce receita de bolis, mais bolis eu num gostis. Aenean aliquam molestie leo, vitae iaculis nisl.</p><p>Nullam volutpat risus nec leo commodo, ut interdum diam laoreet. Sed non consequat odio. Per aumento de cachacis, eu reclamis. Em pé sem cair, deitado sem dormir, sentado sem cochilar e fazendo pose. Delegadis gente finis, bibendum egestas augue arcu ut est.</p><p>Mussum Ipsum, cacilds vidis litro abertis. Quem manda na minha terra sou euzis! Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus. Interessantiss quisso pudia ce receita de bolis, mais bolis eu num gostis. Aenean aliquam molestie leo, vitae iaculis nisl.</p><p>Nullam volutpat risus nec leo commodo, ut interdum diam laoreet. Sed non consequat odio. Per aumento de cachacis, eu reclamis. Em pé sem cair, deitado sem dormir, sentado sem cochilar e fazendo pose. Delegadis gente finis, bibendum egestas augue arcu ut est.</p>",
      },
    },
  }
  err = tmpl.Execute( w, pageData )
  if err != nil {
    w.Write( []byte( err.Error() ) )
  }
}

func hello(w mktp.ProxyResponseWriter, r *mktp.ProxyRequest) {

}

func containerListHtml(w mktp.ProxyResponseWriter, r *mktp.ProxyRequest) {
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
  mktp.FuncMap.Add( mktp.ProxyRootConfig.ProxyNotFound )
  mktp.FuncMap.Add( mktp.ProxyRootConfig.ProxyError )
  mktp.FuncMap.Add( blogNovo )
  mktp.FuncMap.Add( containerListHtml )
  mktp.FuncMap.Add( mktp.ProxyRootConfig.RouteAdd )
  mktp.FuncMap.Add( mktp.ProxyRootConfig.RouteDelete )
  mktp.FuncMap.Add( mktp.ProxyRootConfig.ProxyStatistics )

  mktp.ProxyRootConfig = mktp.ProxyConfig{
    ListenAndServe: ":8888",
  }
  mktp.ProxyRootConfig.RouteAddJs(
    mktp.ProxyRoute{
    // docker run -d --name ghost-blog-demo -p 2368:2368 ghost
      Name: "blog",
      Domain: mktp.ProxyDomain{
        Host: "blog.localhost:8888",
      },
        ProxyEnable: true,
        ProxyServers: []mktp.ProxyUrl{
        {
          Name: "docker 1 - ok",
          Url: "http://localhost:2368",
        },
        {
          Name: "docker 2 - error",
          Url: "http://localhost:2367",
        },
        {
          Name: "docker 3 - error",
          Url: "http://localhost:2367",
        },
      },
    },
  )
  //mktp.ProxyRootConfig = mktp.ProxyConfig{
  //  ListenAndServe: ":8888",
  //  Routes: []mktp.ProxyRoute{
  //    {
  //      // docker run -d --name ghost-blog-demo -p 2368:2368 ghost
  //      Name: "blog",
  //      Domain: mktp.ProxyDomain{
  //        Host: "blog.localhost:8888",
  //      },
  //      ProxyEnable: true,
  //      ProxyServers: []mktp.ProxyUrl{
  //        {
  //          Name: "docker 1 - ok",
  //          Url: "http://localhost:2368",
  //        },
  //        {
  //          Name: "docker 2 - error",
  //          Url: "http://localhost:2367",
  //        },
  //        {
  //          Name: "docker 3 - error",
  //          Url: "http://localhost:2367",
  //        },
  //      },
  //    },
  //    {
  //      Name: "blog novo",
  //      Domain: mktp.ProxyDomain{
  //        NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
  //        ErrorHandle: mktp.ProxyRootConfig.ProxyError,
  //        Host: "go.localhost:8888",
  //      },
  //      ProxyEnable: false,
  //      Handle: mktp.ProxyHandle{
  //        Handle: blogNovo,
  //      },
  //    },
  //    {
  //      Name: "container_list_html",
  //      Domain: mktp.ProxyDomain{
  //        NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
  //        ErrorHandle: mktp.ProxyRootConfig.ProxyError,
  //        Host: "panel.localhost:8888",
  //      },
  //      Path: mktp.ProxyPath{
  //        Path : "/htmlContainerList",
  //        Method: "GET",
  //      },
  //      ProxyEnable: false,
  //      Handle: mktp.ProxyHandle{
  //        Handle: containerListHtml,
  //      },
  //    },
  //    {
  //      Name: "addTest",
  //      Domain: mktp.ProxyDomain{
  //        NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
  //        ErrorHandle: mktp.ProxyRootConfig.ProxyError,
  //        Host: "localhost:8888",
  //      },
  //      Path: mktp.ProxyPath{
  //        Path : "/add",
  //        Method: "POST",
  //        // ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
  //      },
  //      ProxyEnable: false,
  //      Handle: mktp.ProxyHandle{
  //        Handle: mktp.ProxyRootConfig.RouteAdd,
  //      },
  //    },
  //    {
  //      Name: "removeTest",
  //      Domain: mktp.ProxyDomain{
  //        NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
  //        ErrorHandle: mktp.ProxyRootConfig.ProxyError,
  //        Host: "localhost:8888",
  //      },
  //      Path: mktp.ProxyPath{
  //        Path : "/remove",
  //        Method: "POST",
  //        // ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
  //      },
  //      ProxyEnable: false,
  //      Handle: mktp.ProxyHandle{
  //        Handle: mktp.ProxyRootConfig.RouteDelete,
  //      },
  //    },
  //    {
  //      Name: "panel",
  //      Domain: mktp.ProxyDomain{
  //        NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
  //        ErrorHandle: mktp.ProxyRootConfig.ProxyError,
  //        Host: "root.localhost:8888",
  //      },
  //      Path: mktp.ProxyPath{
  //        Path: "/statistics",
  //        Method: "GET",
  //      },
  //      ProxyEnable: false,
  //      Handle: mktp.ProxyHandle{
  //        Handle: mktp.ProxyRootConfig.ProxyStatistics,
  //      },
  //    },
  //  },
  //}
  mktp.ProxyRootConfig.Prepare()
  /*
  b, e := mktp.ProxyRootConfig.MarshalJSON()
  if e != nil {
    fmt.Printf( "MarshalJSON error: %v", e.Error() )
  }
  _ = b
  e = mktp.ProxyRootConfig.UnmarshalJSON( b )
  if e != nil {
    fmt.Printf( "UnmarshalJSON error: %v", e.Error() )
  }
  */
  //go mktp.ProxyRootConfig.VerifyDisabled()

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer( http.Dir( "static" ) ) ) )
  http.HandleFunc("/", mktp.ProxyFunc)
  http.ListenAndServe(mktp.ProxyRootConfig.ListenAndServe, nil)
}