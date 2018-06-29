<html>
<head>
    <title>FileServer:{{.path}}</title>
    <link rel="shortcut icon" href="/static/favicon.ico"/>

</head>
<style>
    div{
        height: 28px;
        align-items: center;
    }
    .a_inline {
        width: 200px;
        overflow: hidden;
        display: inline-block;
    }

    span {
        width: 150px;
        overflow: hidden;
        display: inline-block;
    }
</style>
<body>
<h1>{{.path}}</h1>

<div><a href="/list/{{urlquery .parent}}">返回上一级</a></div>
{{range .info}}
<div>
    <a class="a_inline" href="/list/{{urlquery (print $.path "/" (GetName .))}}">{{GetName .}}</a>
{{if .IsDir}}
    <span><a href="/download/{{urlquery (print $.path "/" (GetName .))}}">下载压缩包</a></span>
{{else}}
    <span>{{GetSize .}}</span>
{{end}}
    <span>{{GetTime .}}</span>
</div>
{{end}}

</body>
</html>