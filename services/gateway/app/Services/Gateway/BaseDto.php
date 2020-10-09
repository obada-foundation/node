<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use App\Services\Gateway\Validation\Rules\DtoArrayKeys;
use App\Services\Gateway\Validation\Rules\HashLink;
use Illuminate\Support\Facades\Validator;
use Illuminate\Validation\ValidationException;
use Spatie\DataTransferObject\DataTransferObject;

class BaseDto extends DataTransferObject {

    public $metadata;

    public $docLinks;

    public $structuredData;

    protected function validate() {
        foreach (['metadata', 'doc_links', 'structured_data'] as $field) {
            if (! is_array($this->{$field})) {
                throw ValidationException::withMessages([$field => "The $field must be an array"]);
            }
        }

        $data  = [];
        $rules = [];

        if ($this->metadata) {
            $data['metadata']          = $this->metadata;
            $rules['metadata']         = 'array';
            $rules['metadata.*.key']   = 'required|string';
            $rules['metadata.*.value'] = 'present';
            $rules['metadata.*']       = ['array', new DtoArrayKeys(['key', 'value'])];
        }

        if ($this->docLinks) {
            $data['doc_links']             = $this->docLinks;
            $rules['doc_links']            = 'array';
            $rules['doc_links.*.name']     = 'required|string';
            $rules['doc_links.*.hashlink'] = ['required', 'url', new HashLink];
            $rules['doc_links.*']          = ['array', new DtoArrayKeys(['name', 'hashlink'])];
        }

        if ($this->structuredData) {
            $data['structured_data']        = $this->structuredData;
            $rules['structured_data']       = 'array';
            $rules['structured_data.*.key'] = 'required|string';
            $rules['structured_data.*']     = ['array', new DtoArrayKeys(['key', 'value'])];
        }

        Validator::make($data, $rules)->validate();
    }
}
