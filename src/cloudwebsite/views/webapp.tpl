{{template "appheader.tpl" .}}

<script src="/static/js/long.js"></script>
<script src="/static/js/webapp/binarystream.js?{{RandString}}"></script>
<script src="/static/js/webapp/app.js?{{RandString}}"></script>
<script src="/static/js/webapp.js?{{RandString}}"></script>

{{template "webapp_login_form.tpl" .}}
{{template "chatpanel.tpl" .}}
{{template "webapp/modals.tpl" .}}
{{template "webapp/menu.tpl" .}}
{{template "fpanel.tpl" .}}

{{template "appfooter.tpl" .}}