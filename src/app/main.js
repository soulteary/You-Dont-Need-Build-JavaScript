define(function (require) {
  var DateTimePicker = require("./components/module1");
  var router = sanRouter.router;
  router.add({ rule: "/", Component: DateTimePicker, target: "#app" });

  return {
    init: function () {
      router.start();
    },
  };
});
