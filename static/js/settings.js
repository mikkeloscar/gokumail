var item = function (placeHolder, name) {
  var output = "<li>" +
               "<div class=\"input-group\">" +
               "<input type=\"text\" class=\"form-control list\" " +
               "name=\"" + name + "[]\" placeholder=\"" + placeHolder + "\">" +
               "<div class=\"input-group-addon\">" +
               "<a href=\"#remove\" title=\"Remove\" class=\"remove-item\" " +
               "tabindex=\"-1\">" +
               "<span class=\"glyphicon glyphicon-remove\"></span>" +
               "</a>" +
               "</div>" +
               "</div>" +
               "</li>";
  return output;
};

function capitalize(str) {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

$(document).ready(function () {
  $('body').on('click', '.remove-item', function () {
    var li = $(this).parent().parent().parent();
    li.remove();
    return false;
  });

  $('body').on('keyup', '.list', function () {
    var li = $(this).parent().parent();
    var prev_val = li.prev().find('.list').val();
    var val_length = 1;
    if (typeof prev_val != 'undefined') {
      val_length = prev_val.length;
    }

    if (li.is(':last-child') && val_length !== 0 && $(this).val().length !== 0) {
      var ul = li.parent();
      var type = $(this).attr("name");
      type = type.substring(0, type.length - 2); // remove []
      ul.append($(item(capitalize(type), type)));
    }
  });
  // $('.add-item').on('click', function () {
  //   var ul = $(this).parent().prev();
  //   if ($(this).hasClass('from')) {
  //     ul.append($(item("From", "from")));
  //   } else if ($(this).hasClass('to')) {
  //     ul.append($(item("To", "to")));
  //   } else if ($(this).hasClass('blacklist')) {
  //     ul.append($(item("Blacklist", "blacklist")));
  //   }
  //   return false;
  // });
});
