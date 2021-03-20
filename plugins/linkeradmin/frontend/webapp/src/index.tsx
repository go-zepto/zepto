import './index.css';
import reportWebVitals from './reportWebVitals';
import LinkerAdmin from 'linkeradmin';
import { ComponentGenerator, Field, Schema } from 'linkeradmin/types';

const admin = new LinkerAdmin({
  target: 'root',
  schemaPath: '/admin/_schema',
});


const ColorPicker: ComponentGenerator = {
  fieldCompGen: (schema: Schema, field: Field): React.FC => {
    console.log(field, schema);
    return (props) => (
      <div>
        Hello World!
      </div>
    );
  },
  inputCompGen: (schema: Schema, field: Field): React.FC => {
    console.log(field, schema);
    return (props) => (
      <div>
        Hello World!
      </div>
    );
  }
}

admin.registerComponentGenerator("color_picker", ColorPicker);

admin.init();

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
