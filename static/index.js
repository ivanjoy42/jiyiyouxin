BackUrl = "";
BodyXHR = null;
HomeDeep = 0;
HomePage = "home";

$(function () {
    url = GetHash();
    Href(url, 0);

    $.ajaxSetup({
        cache: true,
    });
});

function GetCookie(name) {
    var arr,
        reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
    if ((arr = document.cookie.match(reg))) {
        return arr[2];
    } else {
        return null;
    }
}

function GetHash() {
    if (location.hash == "") {
        url = HomePage;
    } else {
        url = location.hash.slice(1);
    }
    return url;
}

function Href(url, direct) {
    if (!url.match(/.+/)) {
        return;
    }

    if (BodyXHR) {
        BodyXHR.abort();
    }

    if (direct < 0 && HomeDeep > 0) {
        BackUrl = url;
        history.go(direct);
        return;
    } else if (direct == undefined && HomeDeep > 0) {
        BackUrl = url;
        history.go(-HomeDeep);
        HomeDeep = 1;
        return;
    } else {
        BodyXHR = $.ajax({
            url: url,
            success: function (res) {
                UpdateHistory(url, direct);
                $("body").html(res);
                scrollTo(0, 0);
            },
        });
    }
}

function UpdateHistory(url, direct) {
    if (direct == 1) {
        history.pushState(url, null, "#" + url);
        HomeDeep++;
    } else if (direct == 0) {
        history.replaceState(url, null, "#" + url);
    }
}

onpopstate = function (event) {
    if (event.state == null) {
        url = GetHash();
        HomeDeep++;
    } else if (BackUrl != "") {
        url = BackUrl;
        BackUrl = "";
        HomeDeep--;
    } else {
        url = event.state;
        HomeDeep--;
    }
    Href(url, 0);
};

function AjaxSubmit(url, direct) {
    var data = {};
    $.each($("form").serializeArray(), function () {
        data[this.name] = this.value;
    });

    $.ajax({
        type: "POST",
        url: url,
        data: data,
        success: function (res) {
            if (direct < 0) {
                history.go(direct);
                HomeDeep = HomeDeep + direct + 1;
            } else {
                $("body").html(res);
                scrollTo(0, 0);
            }
        },
    });
    return false;
}