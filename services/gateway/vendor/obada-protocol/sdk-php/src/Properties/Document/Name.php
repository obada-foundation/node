<?php

declare(strict_types=1);

namespace Obada\Properties\Document;

use Obada\Properties\SimpleProperty;

class Name extends SimpleProperty
{
    public function isValid(): bool
    {
        return is_string($this->value) && strlen(trim($this->value)) > 0;
    }

    public function getValidationMessage(): string
    {
        return 'Name must be not empty string with length more than 0.';
    }
}