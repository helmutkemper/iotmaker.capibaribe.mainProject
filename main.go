package main

import (
  mktp "github.com/helmutkemper/marketPlaceProxy"
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
  Title               string
  Text                string
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
}

type Template struct {
  NavigationMenu      []Link
  ListGroup           []Link
  SidePanel           []SidePanel
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
    Template: Template{
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
          Title: "Sidebar panel widget",
          Text: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
        },
        {
          Title: "Sidebar panel widget",
          Text: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
        },
        {
          Title: "Sidebar panel widget",
          Text: "Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.",
        },
      },
    },
    Content: []Content{
      {
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
  mktp.ProxyRootConfig = mktp.ProxyConfig{
    ListenAndServe: ":8888",
    Routes: []mktp.ProxyRoute{
      {
        // docker run -d --name ghost-blog-demo -p 2368:2368 ghost
        Name: "blog",
        Domain: mktp.ProxyDomain{
          SubDomain: "blog",
          Domain: "localhost",
          Port: "8888",
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
      {
        Name: "blog novo",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "go",
          Domain: "localhost",
          Port: "8888",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: blogNovo,
        },
      },
      {
        Name: "image_list",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path : "/imageList",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ImageWebList,
        },
      },
      {
        Name: "container_list_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path : "/restContainerList",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerWebList,
        },
      },
      {
        Name: "container_stop_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: "^/restContainerStop/(?P<id>[a-fA-F0-9-]{64})$",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerStopById,
        },
      },
      {
        Name: "container_remove_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: "^/restContainerRemove/(?P<id>[a-fA-F0-9-]{64})$",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerRemove,
        },
      },
      {
        Name: "container_kill_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: "^/restContainerKill/(?P<id>[a-fA-F0-9-]{64})$",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerRemove,
        },
      },
      {
        Name: "container_stats_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: "^/restContainerStats/(?P<id>[a-fA-F0-9-]{64})$",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerStatsById,
        },
      },
      {
        Name: "container_stats_all_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path: "/restContainerStatsAll",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerWebStatsLog,
        },
      },
      {
        Name: "container_start_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: "^/restContainerStart/(?P<id>[a-fA-F0-9-]{64})$",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerStart,
        },
      },
      {
        Name: "container_stats_logById",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: "^/restContainerStatsLog/(?P<id>[a-fA-F0-9-]{64})$",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerWebStatsLogById,
        },
      },
      {
        Name: "container_list_html",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path : "/htmlContainerList",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: containerListHtml,
        },
      },
      {
        Name: "container_log",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: `^/containerLog/(?P<id>[a-fA-F0-9-]{64})$`,
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerLogsById,
        },
      },







      {
        Name: "addTest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path : "/add",
          Method: "POST",
          // ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ProxyRootConfig.RouteAdd,
        },
      },
      {
        Name: "removeTest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path : "/remove",
          Method: "POST",
          // ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ProxyRootConfig.RouteDelete,
        },
      },
      {
        Name: "panel",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "root",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path: "/statistics",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ProxyRootConfig.ProxyStatistics,
        },
      },
    },
  }
  mktp.ProxyRootConfig.Prepare()
  go mktp.ProxyRootConfig.VerifyDisabled()

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer( http.Dir( "static" ) ) ) )
  http.HandleFunc("/", mktp.ProxyFunc)
  http.ListenAndServe(mktp.ProxyRootConfig.ListenAndServe, nil)
}