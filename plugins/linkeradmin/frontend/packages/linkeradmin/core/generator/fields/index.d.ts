import { ComponentGeneratorFunc } from "../../../types/generators";
interface FieldGenerator {
    fieldCompGen: ComponentGeneratorFunc;
    inputCompGen: ComponentGeneratorFunc;
}
declare type FieldGenerators = {
    [key: string]: FieldGenerator;
};
export declare const DEFAULT_FIELD_GENERATORS: FieldGenerators;
export {};
