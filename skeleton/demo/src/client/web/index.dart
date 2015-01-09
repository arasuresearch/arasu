library {{.Name}};
import 'dart:html';
import 'package:{{.Name}}/controllers/controllers.dart';
import 'package:{{.Name}}/routes/routes.dart';
import 'package:logging/logging.dart';
import 'package:angular/angular.dart';
import 'package:angular/application_factory.dart';

void main() {
  Logger.root
      ..level = Level.FINEST
      ..onRecord.listen((LogRecord r) => print(r.message));
  Routes routes= new Routes();
  routes.bind(RouteInitializerFn, toValue: bind_routes);

  applicationFactory()
      ..addModule(new Controllers())
      ..addModule(routes)
      ..run();

  querySelector("#sample_text_id")
     ..onClick.listen(reverseText);
           
}

void bind_routes(Router router, RouteViewFactory views) {
  controllers_routes(router, views);
  views.configure({
    'default': ngRoute( path: '/index', view: '/views/layouts/application.html'),
    'index': ngRoute(
        defaultRoute: true,
        enter: (RouteEnterEvent e) => router.go('default', {}, startingFrom: router.root.findRoute(''), replace: true)
      )
  });
} 

void reverseText(MouseEvent event) {
  var text = querySelector("#hello").text;
  var buffer = new StringBuffer();
  for (int i = text.length - 1; i >= 0; i--) {
    buffer.write(text[i]);
  }
  querySelector("#hello").text = buffer.toString();
}
