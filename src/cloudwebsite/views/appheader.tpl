<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>AdminLTE 2 | Starter</title>
  <!-- Tell the browser to be responsive to screen width -->
  <meta content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" name="viewport">
  <link rel="stylesheet" href="static/bower_components/bootstrap/dist/css/bootstrap.min.css">
  <!-- Font Awesome -->
  <link rel="stylesheet" href="static/bower_components/font-awesome/css/font-awesome.min.css">
  <!-- Ionicons -->
  <link rel="stylesheet" href="static/bower_components/Ionicons/css/ionicons.min.css">
  <!-- Theme style -->
  <link rel="stylesheet" href="static/dist/css/AdminLTE.min.css">
  <!-- AdminLTE Skins. We have chosen the skin-blue for this starter
        page. However, you can choose any other skin. Make sure you
        apply the skin class to the body tag so the changes take effect. -->
  <link rel="stylesheet" href="static/dist/css/skins/skin-blue.min.css">

  <link href="https://cdn.bootcss.com/jqueryui/1.12.1/jquery-ui.min.css" rel="stylesheet">
</head>
<body style="background: #d2d6de;">

<!-- REQUIRED JS SCRIPTS -->

<!-- jQuery 3 -->
<script src="https://cdn.bootcss.com/jquery/3.3.1/jquery.min.js"></script>
<script src="https://cdn.bootcss.com/jqueryui/1.12.1/jquery-ui.min.js"></script>
<script src="https://cdn.bootcss.com/jQuery-slimScroll/1.3.8/jquery.slimscroll.min.js"></script>
<script src="https://cdn.bootcss.com/dragscroll/0.0.8/dragscroll.min.js"></script>
<script src="https://cdn.bootcss.com/country-select-js/2.0.1/js/countrySelect.min.js"></script>

<!-- Bootstrap 3.3.7 -->
<script src="static/bower_components/bootstrap/dist/js/bootstrap.min.js"></script>
<!-- AdminLTE App -->
<script src="static/dist/js/adminlte.min.js"></script>

<header class="navbar navbar-expand-lg navbar-light" style="background-color: #e3f2fd;">
  <a class="navbar-brand mr-md-auto" href="/">
    WebSite
  </a>

  <div class="pull-right">
    <a class="mr-2" href="/webapp">
      [app login]
    </a>
    
    {{if .account}}
    <div class="mr-md-2">
      欢迎 <a>{{str2html .account}}</a>
      <a href="/user/logout?{{RandString}}">退出登录</a>
    </div>
    {{end}}
  </div>
</header>
<div class="container-fluid">
  <div class="row">