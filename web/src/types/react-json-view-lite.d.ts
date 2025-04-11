declare module 'react-json-view-lite' {
  import React from 'react';

  export interface JsonViewProps {
    data: any;
    style?: React.CSSProperties;
    className?: string;
    shouldExpandNode?: (
      keyPath: (string | number)[],
      data: any,
      level: number
    ) => boolean;
    indentWidth?: number;
    enableClipboard?: boolean;
    keyRenderer?: (
      key: string | number,
      keyPath: (string | number)[],
      data: any
    ) => React.ReactNode;
    valueRenderer?: (
      value: any,
      keyPath: (string | number)[],
      data: any
    ) => React.ReactNode;
  }

  export const JsonView: React.FC<JsonViewProps>;
} 