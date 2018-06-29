<html>
<head>
    <title>FileServer</title>
</head>
<body>
<div>
    <form method="get" action="/index">
        <label>上下文路径:
            <input type="text" name="contextPath" value="{{.path}}">
        </label>
        <input type="submit" value="提交">
    </form>
</div>
</body>
</html>