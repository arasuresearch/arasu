part of {{.App.Name}}.routes;

void {{.Name}}s_routes(Router router, RouteViewFactory views) {
  views.configure({
    '{{.Name}}s': ngRoute(
      path: '/{{.Name}}s',
      view: '/views/layouts/application.html',
      //view: '/views/layouts/{{.Name}}s.html',
      mount: {
        '': ngRoute(
          path: '/',
          defaultRoute: true,
          view: '/views/{{.Name}}s/index.html'),
        'new': ngRoute(
          path: '/new',
          view: '/views/{{.Name}}s/new.html'),
        'show': ngRoute(
          path: '/:id/show',
          view: '/views/{{.Name}}s/show.html'),
        'edit': ngRoute(
          path: '/:id/edit',
          view: '/views/{{.Name}}s/edit.html')
      })
  });
}  
