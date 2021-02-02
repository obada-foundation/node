<?php

declare(strict_types=1);

namespace Obada\Properties;

class Status extends SimpleProperty
{

    public const FUNCTIONAL        = 'FUNCTIONAL';
    public const NON_FUNCTIONAL    = 'NON_FUNCTIONAL';
    public const DISPOSED          = 'DISPOSED';
    public const STOLEN            = 'STOLEN';
    public const DISABLED_BY_OWNER = 'DISABLED_BY_OWNER';

    public const STATUSES = [
        self::FUNCTIONAL,
        self::NON_FUNCTIONAL,
        self::STOLEN,
        self::DISABLED_BY_OWNER
    ];

    /**
     * @return bool
     */
    public function isValid(): bool
    {
        if (strlen($this->value) > 0) {
            return in_array(strtoupper($this->value), self::STATUSES);
        }

        return true;
    }

    /**
     * @return string
     */
    public function getValidationMessage(): string
    {
        return 'Obit status is not supported. Supported statuses: ' . implode(', ', self::STATUSES);
    }
}
