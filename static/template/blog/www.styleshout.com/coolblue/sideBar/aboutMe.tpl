{{ define "aboutMe" }}
  <div class='{{.Class}}'>
    <h3>{{.Title}}</h3>
    <p>
      <!--a href='index.html'-->
        <img src='{{.Image.Link}}' width='{{.Image.Width}}' height='{{.Image.Height}}' alt='{{.Image.Alt}}' class='align-left' />
      <!--/a-->
      {{ htmlSafe .Text }}
    </p>
  </div>
{{ end }}