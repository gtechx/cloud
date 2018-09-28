{{template "header.tpl" .}}

<div style="position:absolute;top:30%;left:50%; transform:translate(-50%, -50%);">  

{{if .post}}
    {{if not .error}}
        <div class="alert alert-success">恭喜，注册成功</div>
    {{else}}
        <div class="alert alert-danger">注册失败：{{str2html .error}}</div>
        {{template "register_form.tpl" .}}
    {{end}}
{{else}}
    {{template "register_form.tpl" .}}
{{end}}

<br/>
<a href="login">点击这里登录</a>
</div>

{{template "footer.tpl" .}}