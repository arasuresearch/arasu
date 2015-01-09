library {{.App.Name}}.{{.Name}}s;
import 'package:{{.App.Name}}/controllers/controllers.dart';
import 'package:{{.App.Name}}/routes/routes.dart';
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
}

void bind_routes(Router router, RouteViewFactory views) {
  controllers_routes(router, views);
  views.configure({
    'index': ngRoute(
        defaultRoute: true,
        enter: (RouteEnterEvent e) => router.go( '{{.Name}}s', {}, startingFrom: router.root.findRoute(''), replace: true)
      )
  });
} 
