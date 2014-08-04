import 'package:arasu/resource/resource.dart';
import 'package:{{.App.Name}}/components/{{.Name}}s/all.dart';
import 'package:angular/angular.dart';
import 'package:logging/logging.dart';
import 'package:angular/application_factory.dart';

void main() {
  Logger.root
      ..level = Level.FINEST
      ..onRecord.listen((LogRecord r) => print(r.message));
  Module m = new Module();
  m.bind(Index);    
  m.bind(New);  
  m.bind(Edit); 
  m.bind(Show);
  m.bind(RouteInitializerFn, toValue: initRoutes);
  m.bind(NgRoutingUsePushState, toFactory: (_) => new NgRoutingUsePushState.value(false));
  m.bind(Db, toFactory: (v) {
    Db db = new Db();
    db.http = (v as Injector).get(Http);
    return db;
  });
  applicationFactory().addModule(m).run();
}

void initRoutes(Router router, RouteViewFactory view) {
  router.root
    ..addRoute(
      name: 'index',
      path: '/',
      defaultRoute: true,
      enter: view('/views/{{.Name}}s/index.html'))
    ..addRoute(
      name: 'new',
      path: '/new',
      enter: view('/views/{{.Name}}s/new.html'))
    ..addRoute(
      name: 'show',
      path: '/:id/show',
      enter: view('/views/{{.Name}}s/show.html'))
    ..addRoute(
      name: 'edit',
      path: '/:id/edit',
      enter: view('/views/{{.Name}}s/edit.html'));
}  
