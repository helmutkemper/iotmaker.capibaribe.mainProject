{{ define "advertsMenu" }}
<div class="{{.Class}}">

  <h3>{{.Title}}</h3>

  <ul>
    {{ range .Items }}
    <li><a href="{{.Link}}" title="{{.Title}}">{{.Label}}
      <span>{{htmlSafe .Text}}</span></a>
    </li>
    {{ end }}
  </ul>

</div>
{{ end }}