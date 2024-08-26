import React, { useState } from 'react';
import "../Stylesheets/navegador.css"

export default function Comandos({newIp="localhost"}){
    const [textValue, setTextValue] = useState('');
    const [textExit, setTextExit] = useState('');

    const handleTextChange = (event) => {
        setTextValue(event.target.value);
    };

    const sendData = async (e) => {
        e.preventDefault();
        const data = {
            text: textValue
        };
        
        try {
            const response = await fetch(`http://localhost:8080/analizar`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });
    
            if (!response.ok) {
                throw new Error('Error al enviar datos');
            }
    
            const responseData = await response.json();
            console.log('Respuesta del servidor:', responseData);
            console.log('Respuesta del metodo ',responseData.message)
            setTextExit(responseData.message)
           
        } catch (error) {
            console.error('Error:', error);
        }

    }

    return(
        <div className='contenedorEjecutar'>
            <div id="espacio">&nbsp;&nbsp;&nbsp;</div>
            <table>
                <tbody>

                    <tr><td><p><strong>ENTRADA</strong></p></td></tr>

                    <tr>
                        <td>
                            <textarea
                                className='entrada'
                                value={textValue}
                                onChange={handleTextChange}
                                placeholder='Ingrese comandos'
                                id='inputComands'
                            />
                        </td>
                    </tr>

                    <tr><td><strong><p>SALIDA</p></strong></td></tr>

                    <tr>
                        <td>
                            <textarea
                                className='entrada'
                                value={textExit}
                                id='inputComands'
                            />
                        </td>
                    </tr>

                    <tr>
                        <td style={{textAlign:'center'}}>
                            <button type="button" className="btn btn-primary" onClick={(e) => sendData(e)}>Ejecutar</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    );
}