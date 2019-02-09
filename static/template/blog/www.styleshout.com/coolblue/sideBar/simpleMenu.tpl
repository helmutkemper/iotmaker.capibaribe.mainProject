{{ define "simpleMenu" }}
    <div class="{{.Class}}">

      <h3>{{.Title}}</h3>
      <ul>
        {{ range .Items }}<li><a href="{{.Link}}" title="{{.Title}}">{{.Label}}</a></li>{{ end }}
      </ul>

    </div>
{{ end }}