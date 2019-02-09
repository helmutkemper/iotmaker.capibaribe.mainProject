{{ define "postContent" }}
    {{if ne .SimpleTemplate true}}
        {{range .Post}}
            <article class="post">
                <div class="primary">
                    <h2>
                        {{.Title}}
                    </h2>
                    {{$element := .LabelKeyWords}}
                    {{if ne .Label ""}}
                        <p class="post-info">
                            <span>
                                {{.Label}}
                            </span> {{range $k, $v := .LabelKeyWords}}
                            <a href="{{$v.Link}}" title="{{$v.Title}}" >
                                {{$v.Label}}
                            </a>{{if nLast $k $element}},{{end}} {{end}}
                        </p>
                    {{end}}
                    <div class="image-section">
                        <img src="{{.Image.Link}}" alt="{{.Image.Alt}}" height="{{.Image.Height}}" width="{{.Image.Width}}"/>
                    </div>
                    {{htmlSafe .Text}}
                </div>
                <aside>
                    <p class="dateinfo">
                        {{.PostInfo.Month}}
                        <span>
                            {{.PostInfo.Day}}
                        </span>
                    </p>
                    <div class="post-meta">
                        <h4>
                            {{$.PostInfo}}
                        </h4>
                        <ul>
                            <li class="user">
                                <a href="#">
                                    {{.PostInfo.User}}
                                </a>
                            </li>
                            <li class="time">
                                <a href="#">
                                    {{.PostInfo.Time}}
                                </a>
                            </li>
                            <li class="comment">
                                <a href="#">
                                    {{.PostInfo.Comments}} {{$.Comments}}
                                </a>
                            </li>
                            <li class="permalink">
                                <a href="{{.PostInfo.Permalink}}">
                                    Permalink
                                </a>
                            </li>
                        </ul>
                    </div>
                </aside>
            </article>
        {{end}}
    {{else}}
        {{range .Post}}
            <div class="main-content">
                <h2>
                    {{.Title}}
                </h2>
                {{$element := .LabelKeyWords}}
                {{if ne .Label ""}}
                    <p class="post-info">
                        <span>
                            {{.Label}}
                        </span> {{range $k, $v := .LabelKeyWords}}
                        <a href="{{$v.Link}}">
                            {{$v.Label}}
                        </a>{{if nLast $k $element}},{{end}} {{end}}
                    </p>{{end}}
                <p>
                    <a href="#">
                        <img src="{{.Image.Link}}" alt="{{.Image.Alt}}" height="{{.Image.Height}}" width="{{.Image.Width}}"  class="align-left"/>
                    </a>
                    {{htmlSafe .Text}}
                </p>
            </div>
        {{end}}
    {{end}}
{{end}}