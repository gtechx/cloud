{{template "header.tpl" .}}

<div class="list-group" style="position:absolute;top:50%;left:50%; transform:translate(-50%, -0%);">
  <a href="/user/app/index" class="list-group-item list-group-item-action text-center">
    我的应用
  </a>
  <a href="/user/appdata/index" class="list-group-item list-group-item-action text-center">
    我的应用数据
  </a>

  <a href="/admin/account/index" class="list-group-item list-group-item-action text-center mt-4">
    用户管理
  </a>
  <a href="/admin/appdata/index" class="list-group-item list-group-item-action text-center">
    应用管理
  </a>
  <a href="data" class="list-group-item list-group-item-action text-center">
    数据管理
  </a>
  <a href="data" class="list-group-item list-group-item-action text-center">
    在线玩家管理
  </a>
</div>

{{template "footer.tpl" .}}