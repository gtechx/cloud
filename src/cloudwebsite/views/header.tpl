<!doctype html>
<html lang="en">
<head>
<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<title>register</title>
<meta content="GTech Inc." name="Copyright" />
<link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/bootstrap-table/1.12.1/bootstrap-table.min.css">
<link rel="stylesheet" href="/static/css/open-iconic-bootstrap.min.css">
<link rel="stylesheet" href="https://cdn.bootcss.com/jqueryui/1.12.1/jquery-ui.min.css">
<link href="https://cdn.bootcss.com/country-select-js/2.0.1/css/countrySelect.min.css" rel="stylesheet">

<link href="https://cdn.bootcss.com/animate.css/3.5.2/animate.min.css" rel="stylesheet">
<link href="https://cdn.bootcss.com/loaders.css/0.1.2/loaders.min.css" rel="stylesheet">
<script src="/static/js/md5.min.js"></script>
</head>
<body>

<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js"></script>
<script src="https://cdn.bootcss.com/bootstrap/4.0.0/js/bootstrap.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/bootstrap-table/1.12.1/bootstrap-table.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/bootstrap-table/1.12.1/locale/bootstrap-table-zh-CN.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/jqueryui/1.12.1/jquery-ui.min.js"></script>
<script src="https://cdn.bootcss.com/country-select-js/2.0.1/js/countrySelect.min.js"></script>


<header class="navbar navbar-expand-lg navbar-light" style="background-color: #e3f2fd;">
  <a class="navbar-brand mr-md-auto" href="/">
    WebSite
  </a>
  <a class="mr-2" href="/webapp">
    [app login]
  </a>
  
  {{if .account}}
  <div class="mr-md-2">
    欢迎 <a>{{str2html .account}}</a>
    <a href="/user/logout?{{RandString}}">退出登录</a>
  </div>
  {{end}}
</header>
<div class="container-fluid">
  <div class="row">
  {{if .account}}
    <nav class="col-md-2 bg-light border-right">
      <div class="sidebar-sticky">
        <ul class="nav nav-pills flex-column">
         <li class="nav-item">
            <a href="/user/index" class="nav-link {{if compare .nav "user"}}active{{end}}">
            综合
            </a>
          </li>
          <li class="nav-item">
            <a href="/user/app/index" class="nav-link {{if compare .nav "userapp"}}active{{end}}">
            应用
            </a>
          </li>
          <li class="nav-item">
            <a href="/user/appdata/index" class="nav-link {{if compare .nav "userappdata"}}active{{end}}">
            应用数据
            </a>
          </li>
          <li class="nav-item">
            <a href="/user/online/index" class="nav-link {{if compare .nav "useronline"}}active{{end}}">
            应用在线玩家
            </a>
          </li>
          <li class="nav-item">
            <a href="/user/myappdata/index" class="nav-link {{if compare .nav "usermyappdata"}}active{{end}}">
            我的应用数据
            </a>
          </li>

          <li class="nav-item mt-4">
          </li>
          {{if .priv.Adminadmin}}
          <li class="nav-item">
            <a href="/admin/admin/index" class="nav-link {{if compare .nav "adminadmin"}}active{{end}}">
            管理员管理
            </a>
          </li>
          {{end}}
          {{if .priv.Adminaccount}}
          <li class="nav-item">
            <a href="/admin/account/index" class="nav-link {{if compare .nav "adminaccount"}}active{{end}}">
            账号管理
            </a>
          </li>
          {{end}}
          {{if .priv.Adminapp}}
          <li class="nav-item">
            <a href="/admin/app/index" class="nav-link {{if compare .nav "adminapp"}}active{{end}}">
            应用管理
            </a>
          </li>
          {{end}}
          {{if .priv.Adminappdata}}
          <li class="nav-item">
            <a href="/admin/appdata/index" class="nav-link {{if compare .nav "adminappdata"}}active{{end}}">
            应用数据管理
            </a>
          </li>
          {{end}}
          {{if .priv.Adminonline}}
          <li class="nav-item">
            <a href="data" class="nav-link">
            在线玩家管理
            </a>
          </li>
          {{end}}
        </ul>
        </div>
    </nav>
    <main role="main" class="col-md-9 col-lg-10 px-0">
    <div class="p-3 bg-light">
    {{end}}