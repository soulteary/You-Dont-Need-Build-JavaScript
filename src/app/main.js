define(["./components/container"], function (Container, require) {
  // var Container = require("./components/container");
  var router = sanRouter.router;
  router.add({ rule: "/", Component: Container, target: "#app" });

  return {
    init: function () {
      router.start();
    },
  };
});
