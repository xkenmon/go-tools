<html>
<head>
    <title>FileServer:{{.path}}</title>
    <link rel="shortcut icon" href="/static/favicon.ico"/>

</head>
<style>
    a {
        width: 300px;
        overflow: hidden;
    }

    span {
        width: 300px;
        overflow: hidden;
    }
</style>
<body>
<h1>{{.path}}</h1>

<div><a href="/list/{{.parent}}">返回上一级</a></div>
{{range .info}}
<div>
    <a href="/list/{{$.path}}/{{GetName .}}">{{GetName .}}</a> <span>{{GetTime .}}</span>
{{if .IsDir}}
    <a href="/download/{{$.path}}/{{GetName .}}">下载压缩包</a>
{{else}}
    <span>{{GetSize .}}</span>
{{end}}
</div>
{{end}}

</body>
</html>