
function loadVerifyCode() {
    // $.get("captcha?type=id",
    //     function (data) {
    //         if (!data.Ok) {
    //             alert(data.Err);
    //             $('#info').html(data.Err);
    //             return;
    //         }
    //         $('#codeImg').attr("src", "/captcha?type=pic&name=" + data.Id + ".png");
    //         $('#imgId').val(data.Id);
    //     }, "json");

    fetch("captcha?type=id", {
        method: "GET"
    }).then(function (res) {
        if (res.ok) {
            res.json().then(function (data) {
                if (!data.Ok) {
                    return;
                }
                $('#codeImg').attr("src", "/captcha?type=pic&name=" + data.Id + ".png");
                $('#imgId').val(data.Id);
            });
        }
    }).catch(function (err) {
        console.log(err);
    });
}

function showNotify(nType, title, msg) {
    var mIcon = 'glyphicon glyphicon-info-sign';
    switch (nType) {
        case 'success':
            mIcon = 'glyphicon glyphicon-ok-circle';
            break;
        case 'info':
            mIcon = 'glyphicon glyphicon-info-sign';
            break;
        case 'warning':
            mIcon = 'glyphicon glyphicon-warning-sign';
            break;
        case 'danger':
            mIcon = 'glyphicon glyphicon-remove-circle';
            break;
    }
    $.notify({
        icon: mIcon,
        title: title,
        message: msg,
        // url: '',
        // target: '_blank'
    }, {
            element: 'body',
            newest_on_top: true,
            type: nType,
            position: null,
            allow_dismiss: true,
            placement: {
                from: "top",
                align: "right"
            },
            offset: {
                x: 0,
                y: 50
            },
            animate: {
                enter: 'animated fadeInDown',
                exit: 'animated fadeOutUp'
            },
            onShow: null,
            onShown: null,
            onClose: null,
            onClosed: null,
        });
}