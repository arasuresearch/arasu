library {{.App.Name}}.controllers.{{.Name}}s;
import 'package:{{.App.Name}}/models/models.dart' show {{.Cname}};
import 'package:arasu/resource/resource.dart';
import 'package:angular/angular.dart';
import 'package:logging/logging.dart';

part 'new.dart';
part 'edit.dart';
part 'show.dart';
part 'index.dart';

void bind_{{.Name}}s_actions(Module m){
  m.bind({{.Cname}}sIndex);
  m.bind({{.Cname}}sShow);
  m.bind({{.Cname}}sNew);
  m.bind({{.Cname}}sEdit);  
}
