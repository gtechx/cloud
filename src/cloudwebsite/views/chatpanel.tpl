<div class="col-md-3 hide" id="chatpanel" style="position:absolute;">
    <!-- DIRECT CHAT PRIMARY -->
    <div class="box box-primary direct-chat direct-chat-primary">
        <div class="box-header with-border">
            <h2 class="box-id hide">Direct Chat</h2>
            <i class="fa fa-group hide"><h3 class="box-title"></i>Direct Chat</h3>

            <div class="box-tools pull-right">
                <span data-toggle="tooltip" title="" class="badge bg-light-blue hide" data-original-title="3 New Messages">3</span>
                <button type="button" class="btn btn-box-tool" data-widget="collapse"><i class="fa fa-minus"></i>
                </button>
                <button type="button" class="btn btn-box-tool roomuser hide" data-toggle="tooltip" title="" data-widget="chat-pane-toggle" data-original-title="Contacts">
                    <i class="fa fa-comments"></i>
                </button>
                <button type="button" class="btn btn-box-tool" data-widget="remove"><i class="fa fa-times"></i></button>
            </div>
        </div>
        <!-- /.box-header -->
        <div class="box-body">
            <!-- Conversations are loaded here -->
            <div class="direct-chat-messages">
            </div>
            <!--/.direct-chat-messages-->

            <!-- Contacts are loaded here -->
            <div class="direct-chat-contacts">
                <ul class="contacts-list">
                    <li id="contacts-list-item" class="hide">
                        <a href="#">
                            <img class="contacts-list-img" src="static/dist/img/user1-128x128.jpg" alt="User Image">

                            <div class="contacts-list-info">
                                <span class="contacts-list-name">
                                    Count Dracula
                                    <small class="contacts-list-date pull-right">2/28/2015</small>
                                </span>
                            <span class="contacts-list-msg">How have you been? I was...</span>
                            </div>
                            <!-- /.contacts-list-info -->
                        </a>
                    </li>
                    <!-- End Contact Item -->
                </ul>
            <!-- /.contatcts-list -->
            </div>
        </div>
        <!-- /.box-body -->
        <div class="box-footer">
            <div class="input-group">
                <input type="text" name="message" placeholder="Type Message ..." class="form-control" />
                    <span class="input-group-btn">
                    <button type="submit" class="btn btn-primary btn-flat">Send</button>
                    </span>
            </div>
        </div>
        <!-- /.box-footer-->
    </div>
    <!--/.direct-chat -->
</div>

<script>
$( function() {
    // $( "#chatpanel" ).draggable({handle: "#chatpanelheader", cursor: "move"});
    // $("#chatpanelpart").resizable({handles: "se", minWidth: 250, maxWidth:500, minHeight:500, maxHeight:650});
    // $("#chatpanelpart").css("height", 500).css("width", 250);
    $(".direct-chat .box-header h3").html("Title");
    $( "#chatpanel" ).draggable({handle: ".direct-chat .box-header", cursor: "move"});
} );

var chatpanellist = {};
var chatpanelroomlist = {};
var curz = 10;
function openChatPanel(data) {
    //onclick="$('#chatpanel').addClass('hide');"
    if(chatpanellist["chatpanel-" + data.nickname] != null){
        chatpanellist["chatpanel-" + data.nickname].css("z-index", curz);
        curz++;
        return chatpanellist["chatpanel-" + data.nickname];
    }
    var newchatpanel = $( "#chatpanel" ).clone();
    newchatpanel.find(".direct-chat").directChat();
    newchatpanel.find(".direct-chat").boxWidget();
    newchatpanel.attr("id", "chatpanel-" + data.nickname);
    newchatpanel.removeClass("hide");
    if(data.comment != ""){
        newchatpanel.find(".direct-chat .box-header h3").html(data.comment + "(" + data.nickname + ")");
    } else {
        newchatpanel.find(".direct-chat .box-header h3").html(data.nickname);
    }
    newchatpanel.draggable({handle: ".direct-chat .box-header", cursor: "move"});
    newchatpanel.css("z-index", curz);
    curz++;
    chatpanellist["chatpanel-" + data.nickname] = newchatpanel;

    newchatpanel.mousedown(function(){
        newchatpanel.css("z-index", curz);
        curz++;
    });

    newchatpanel.find(".box-footer button").click(function(){
        var msginput = newchatpanel.find(".box-footer input");
        startSendMessage(data, msginput.val());
        msginput.val('');
    });

    newchatpanel.on('removed.boxwidget', function (event) {
      newchatpanel.remove();
      delete chatpanellist["chatpanel-" + data.nickname];
      return true;
    });

    $( "#chatpanel" ).parent().append(newchatpanel);
    return newchatpanel;
    // $( "#chatpanel" ).removeClass("hide");
    // $(".direct-chat .box-header h3").html(data.nickname);
}

function openRoomChatPanel(data) {
    //onclick="$('#chatpanel').addClass('hide');"
    if(chatpanelroomlist["chatpanel-room-" + data.rid] != null){
        chatpanelroomlist["chatpanel-room-" + data.rid].css("z-index", curz);
        curz++;
        return chatpanelroomlist["chatpanel-room-" + data.rid];
    }
    var newchatpanel = $( "#chatpanel" ).clone();
    newchatpanel.find(".direct-chat").directChat();
    newchatpanel.find(".direct-chat").boxWidget();
    newchatpanel.attr("id", "chatpanel-room-" + data.rid);
    newchatpanel.removeClass("hide");
    newchatpanel.find(".direct-chat .box-header h3").html(data.roomname);

    newchatpanel.find(".direct-chat .box-header .fa").removeClass("hide");
    newchatpanel.find(".roomuser").removeClass("hide");

    newchatpanel.draggable({handle: ".direct-chat .box-header", cursor: "move"});
    newchatpanel.css("z-index", curz);
    curz++;
    chatpanelroomlist["chatpanel-room-" + data.rid] = newchatpanel;

    newchatpanel.mousedown(function(){
        newchatpanel.css("z-index", curz);
        curz++;
    });

    newchatpanel.find(".box-footer button").click(function(){
        var msginput = newchatpanel.find(".box-footer input");
        startSendRoomMessage(data, msginput.val());
        msginput.val('');
    });

    newchatpanel.on('removed.boxwidget', function (event) {
      newchatpanel.remove();
      delete chatpanelroomlist["chatpanel-room-" + data.rid];
      return true;
    });

    $( "#chatpanel" ).parent().append(newchatpanel);
    reqRoomUserList(data.rid);
    return newchatpanel;
    // $( "#chatpanel" ).removeClass("hide");
    // $(".direct-chat .box-header h3").html(data.nickname);
}

function getChatTitle() {
    return $(".direct-chat .box-header h3").html();
}

function getChatId() {
    return $(".direct-chat .box-header h2").html();
}

function addMessage(msg) {
    var chatpanel = chatpanellist["chatpanel-" + msg.nickname];
    if(chatpanel == null){
        chatpanel = openChatPanel(frienddatabyid[msg.from]);
    }
    var html = '<div class="direct-chat-msg">' +
        '<div class="direct-chat-info clearfix">';
            if(frienddatabyid[msg.from] && frienddatabyid[msg.from].comment != "")
                html += '<span class="direct-chat-name pull-left">' + frienddatabyid[msg.from].comment + '(' + msg.nickname + ')' + '</span>';
            else
                html += '<span class="direct-chat-name pull-left">' + msg.nickname + '</span>';
            html += '<span class="direct-chat-timestamp pull-right">' + new Date(parseInt(msg.timestamp) * 1000).format() + '</span>' +
        '</div>' +
        '<img class="direct-chat-img" src="static/dist/img/user1-128x128.jpg" alt="Message User Image">' +
        '<div class="direct-chat-text">'+
            msg.message + 
        '</div>'+
    '</div>';

    chatpanel.find(".direct-chat-messages").append(html);
    chatpanel.find(".direct-chat-messages").scrollTop(9999);
}

function addSendMessage(msg) {
    var chatpanel = chatpanellist["chatpanel-" + msg.nickname];
    if(chatpanel == null){
        chatpanel = openChatPanel(frienddatabyid[msg.to]);
    }

    var html = '<div class="direct-chat-msg right">' +
        '<div class="direct-chat-info clearfix">' +
            '<span class="direct-chat-name pull-right">' + msg.nickname + '</span>' +
            '<span class="direct-chat-timestamp pull-left">' + new Date().format("yyyy/MM/dd hh:mm:ss") + '</span>' +
        '</div>' +
        '<img class="direct-chat-img" src="static/dist/img/user1-128x128.jpg" alt="Message User Image">' +
        '<div class="direct-chat-text">'+
            msg.message + 
        '</div>'+
    '</div>';

    chatpanel.find(".direct-chat-messages").append(html);

    //console.info($(".direct-chat-msg:last").scrollTop());
    chatpanel.find(".direct-chat-messages").scrollTop(9999);
}
//addMessage({nickname:"WYQ", timestamp:"2018/6/11", message:"Hello, How are you?"});
//addSendMessage({nickname:"WLN", timestamp:"2018/6/11", message:"Hello, I'm fine."});

function startSendMessage(data, message) {
    var msg = {};
    //msg.nickname = userdata.nickname;
    //msg.timestamp = new Date().getTime() / 1000;//new Date().format("yyyy/MM/dd hh:mm:ss");
    msg.message = message;
    msg.to = data.who;

    addSendMessage({nickname:userdata.nickname, to:data.who, message: message});
    sendMessage(msg);
}

function addRoomMessage(msg) {
    var chatpanel = chatpanelroomlist["chatpanel-room-" + msg.rid];
    if(chatpanel == null){
        chatpanel = openRoomChatPanel(roomdata[msg.rid]);
    }
    var html = '<div class="direct-chat-msg">' +
        '<div class="direct-chat-info clearfix">';
            // if(frienddatabyid[msg.who] && frienddatabyid[msg.who].comment != "")
            //     html += '<span class="direct-chat-name pull-left">' + frienddatabyid[msg.who].comment + '(' + msg.nickname + ')' + '</span>';
            // else
                html += '<span class="direct-chat-name pull-left">' + msg.nickname + '</span>';
            html += '<span class="direct-chat-timestamp pull-right">' + new Date(parseInt(msg.timestamp) * 1000).format() + '</span>' +
        '</div>' +
        '<img class="direct-chat-img" src="static/dist/img/user1-128x128.jpg" alt="Message User Image">' +
        '<div class="direct-chat-text">'+
            msg.message + 
        '</div>'+
    '</div>';

    chatpanel.find(".direct-chat-messages").append(html);
    chatpanel.find(".direct-chat-messages").scrollTop(9999);
}

function addRoomSendMessage(msg) {
    var chatpanel = chatpanelroomlist["chatpanel-room-" + msg.rid];
    if(chatpanel == null){
        chatpanel = openRoomChatPanel(roomdata[msg.rid]);
    }
    var html = '<div class="direct-chat-msg right">' +
        '<div class="direct-chat-info clearfix">' +
            '<span class="direct-chat-name pull-right">' + msg.nickname + '</span>' +
            '<span class="direct-chat-timestamp pull-left">' + new Date().format() + '</span>' +
        '</div>' +
        '<img class="direct-chat-img" src="static/dist/img/user1-128x128.jpg" alt="Message User Image">' +
        '<div class="direct-chat-text">'+
            msg.message + 
        '</div>'+
    '</div>';

    chatpanel.find(".direct-chat-messages").append(html);

    //console.info($(".direct-chat-msg:last").scrollTop());
    chatpanel.find(".direct-chat-messages").scrollTop(9999);
}
//addMessage({nickname:"WYQ", timestamp:"2018/6/11", message:"Hello, How are you?"});
//addSendMessage({nickname:"WLN", timestamp:"2018/6/11", message:"Hello, I'm fine."});

function startSendRoomMessage(data, message) {
    // var msg = {};
    // //msg.nickname = userdata.nickname;
    // //msg.timestamp = new Date().getTime() / 1000;//new Date().format("yyyy/MM/dd hh:mm:ss");
    // msg.message = message;
    // //msg.who = data.who;
    // msg.rid = data.rid;

    addRoomSendMessage({rid:data.rid, nickname:userdata.nickname, to:data.nickname, message: message});
    reqSendRoomMessage(data.rid, message);
}

// <li id="contacts-list-item" class="hide">
//     <a href="#">
//         <img class="contacts-list-img" src="static/dist/img/user1-128x128.jpg" alt="User Image">

//         <div class="contacts-list-info">
//             <span class="contacts-list-name">
//                 Count Dracula
//                 <small class="contacts-list-date pull-right">2/28/2015</small>
//             </span>
//         <span class="contacts-list-msg">How have you been? I was...</span>
//         </div>
//         <!-- /.contacts-list-info -->
//     </a>
// </li>
function addRoomUser(user) {
    var newitem = $( "#contacts-list-item" ).clone();

    newitem.mouseup(function(e){
        if(e.button===2){
            showRoomUserListMenu(e, user);
            stopPropagation(e);//调用停止冒泡方法,阻止document方法的执行
        }
    });

    newitem.find(".contacts-list-name").html(user.nickname);
    newitem.attr("id", "contacts-list-item-" + user.dataid);
    newitem.removeClass("hide");

    $( ".direct-chat-contacts .contacts-list" ).append(newitem);
}
</script>