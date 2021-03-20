import { ComponentGenerator } from './types';
interface LinkerAdminOptions {
    schemaPath?: string;
    target: string;
    defaultRowClick?: 'edit' | 'show';
}
export default class LinkerAdmin {
    schemaPath: string;
    target: string;
    defaultRowClick: 'edit' | 'show';
    constructor({ schemaPath, target, defaultRowClick, }: LinkerAdminOptions);
    registerComponentGenerator(name: string, compGen: ComponentGenerator): void;
    init(): Promise<void>;
}
export {};
